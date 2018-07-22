package utilities

import (
	"errors"
	"log"

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

//SetUserForConfirmation sets a user in redis
func SetUserForConfirmation(id string) (string, error) {

	key := uuid.Must(uuid.NewV4()).String()

	_, err := client.Set(key, id, 0).Result()

	if err != nil {
		log.Fatal(err)
		return "", nil
	}

	return key, nil

}

func GetUserForConfirmation(key string) (string, error) {

	id, err := client.Get(key).Result()

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	if id == "" {
		err = errors.New("Please supply a valid key")
		return "", err
	}

	return id, nil
}
