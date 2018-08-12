package utilities

import (
	"fmt"
	"log"
)

//GetSession gets the values of a session
func GetSession(sessionID string) (map[string]string, error) {

	session, err := client.HGetAll(sessionID).Result()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(session) == 0 {
		return nil, fmt.Errorf("Session has expired")
	}

	return session, nil
}
