package comments

import (
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

//Comment is the data structure of a comment on an item
type Comment struct {
	ID         string
	ItemID     string
	Username   string
	Comment    string
	ReplyCount int64
	Replies    []Reply
	CreatedAt  float64
	UpdatedAt  int64
}

//Reply a comment by a user
type Reply struct {
	ID        string
	Username  string
	CommentID string
	Comment   string
	CreatedAt float64
}

//Create a comment for an item
func (comment *Comment) Create() error {
	comment.ID = uuid.Must(uuid.NewV4()).String()
	comment.CreatedAt = float64(time.Now().Unix())
	comment.UpdatedAt = 0

	var z redis.Z
	z.Score = comment.CreatedAt
	z.Member = comment.ID

	pipeline := client.Pipeline()

	pipeline.ZAdd(comment.ItemID, z) // add a comment id to a zset representing a an item

	fields := map[string]interface{}{
		"username":   comment.Username,
		"comment":    comment.Comment,
		"created_at": comment.CreatedAt,
		"updated_at": comment.UpdatedAt,
	}

	pipeline.HMSet(comment.ID, fields)

	_, err := pipeline.Exec()

	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

//Create a reply to a comment by a user
func (reply *Reply) Create() error {
	reply.ID = uuid.Must(uuid.NewV4()).String()
	reply.CreatedAt = float64(time.Now().Unix())

	fields := map[string]interface{}{
		"username":   reply.Username,
		"comment":    reply.Comment,
		"created_at": reply.CreatedAt,
	}

	z := redis.Z{Score: reply.CreatedAt, Member: reply.ID}

	pipeline := client.Pipeline()
	pipeline.ZAdd("replies:"+reply.CommentID, z)
	pipeline.HMSet(reply.ID, fields)
	_, err := pipeline.Exec()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

//Get a reply
func (reply *Reply) Get() error {
	resp, err := client.HGetAll(reply.ID).Result()

	if err != nil {
		log.Println(err)
		return err
	}

	reply.Comment = resp["comment"]
	reply.CreatedAt, _ = strconv.ParseFloat(resp["created_at"], 64)
	reply.Username = resp["username"]

	return nil
}

//Get a comment from database
func (comment *Comment) Get() error {
	resp, err := client.HGetAll(comment.ID).Result()
	if err != nil {
		log.Println(err)
		return err
	}

	replyCount, err := client.ZCount("replies:"+comment.ID, "-inf", "+inf").Result()

	if err != nil {
		log.Println(err)
		return err
	}

	comment.Username = resp["username"]
	comment.Comment = resp["comment"]
	comment.ReplyCount = replyCount

	return nil

}

//GetReplies to a comment
func (comment *Comment) GetReplies() error {

	opt := redis.ZRangeBy{Min: "-inf", Max: "+inf"}

	resp, err := client.ZRangeByScoreWithScores("replies:"+comment.ID, opt).Result()

	if err != nil {
		log.Println(err)
		return err
	}

	for _, v := range resp {
		var reply Reply
		id := v.Member.(string)
		res, err := client.HGetAll(id).Result()
		if err != nil {
			log.Println(err)
			continue
		}

		reply.Username = res["username"]
		reply.Comment = res["comment"]
		reply.ID = id
		reply.CreatedAt = v.Score

		comment.Replies = append(comment.Replies, reply)

	}
	return nil

}

//Delete a comment from the database
func (comment *Comment) Delete() error {
	pipeline := client.Pipeline()

	pipeline.ZRem(comment.ItemID, comment.ID)
	pipeline.Del(comment.ID)
	_, err := pipeline.Exec()

	if err != nil {
		log.Println(err)
		return err
	}

	opt := redis.ZRangeBy{Min: "-inf", Max: "+inf"}
	resp, err := client.ZRangeByScoreWithScores("replies:"+comment.ID, opt).Result()

	if err != nil {
		log.Println(err)
		return err
	}

	for _, v := range resp {
		id := v.Member.(string)
		pipeline.Del(id)
	}

	_, err = pipeline.Exec()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Delete a reply
func (reply *Reply) Delete() error {
	pipeline := client.Pipeline()
	pipeline.ZRem("replies:"+reply.CommentID, reply.ID)
	pipeline.Del(reply.ID)
	_, err := pipeline.Exec()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Update a comment
func (comment *Comment) Update() error {
	updatedAt := time.Now().Unix()
	fields := map[string]interface{}{
		"comment":    comment.Comment,
		"updated_at": updatedAt,
	}
	_, err := client.HMSet(comment.ID, fields).Result()

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Update a reply to a comment
func (reply *Reply) Update() error {
	updatedAt := time.Now().Unix()

	fields := map[string]interface{}{
		"comment":    reply.Comment,
		"updated_at": updatedAt,
	}

	_, err := client.HMSet(reply.ID, fields).Result()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//GetItemComments gets all comments associated with an item
func (comment *Comment) GetItemComments() ([]Comment, error) {
	opts := redis.ZRangeBy{
		Max: "+inf",
		Min: "-inf",
	}
	resp, err := client.ZRangeByScoreWithScores(comment.ItemID, opts).Result()

	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	var comments []Comment

	for _, v := range resp {
		var comment Comment
		id := v.Member.(string)
		comment.ID = id
		comment.CreatedAt = v.Score
		err := comment.Get()

		if err != nil {
			log.Println(err)
			continue
		}

		comments = append(comments, comment)
	}

	return comments, nil

}
