package comments

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // postgres driver
)

//Comment is the data structure of a comment on an item
type Comment struct {
	ID        string
	ItemID    string
	UserID    string
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("postgres", "user=help.ng passowrd=this is the password for help.ng sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

//Create a comment for an item
func (comment *Comment) Create() error {
	query := "INSERT INTO items (item_id, user_id, comment) values ($1, $2, $3,)returning id"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	err = stmt.QueryRow(comment.ItemID, comment.UserID, comment.Comment).Scan(&comment.ID)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return err
}

//Update a comment for an item
func (comment *Comment) Update() error {
	query := "UPDATE comment SET comment = $1 where id = $2"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(comment.Comment)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

//Get a comment from database
func (comment *Comment) Get() error {
	query := "SELECT comment, user_id, created_at, item_id from comments where id = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	err = stmt.QueryRow(comment.ID).Scan(&comment.Comment, &comment.UserID, &comment.CreatedAt, &comment.ItemID)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return err
}

//Delete a comment from the database
func (comment *Comment) Delete() error {
	query := "DELETE FROM comments where id = $id"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(comment.ID)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return err
}
