package utilities

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

//HashPassword hashes a password
func HashPassword(password string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(bytes), err
}

//CheckPassword checks a password if it matches the hash
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
