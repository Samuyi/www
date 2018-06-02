package user

import (
	"database/sql"
	"log"
	"time"

	utilities "github.com/Samuyi/golang-utilities"
	"github.com/Samuyi/www/models/items"
	_ "github.com/lib/pq" // postgres driver
)

//User data structure
type User struct {
	ID        string
	Name      string
	Email     string
	Ratings   string
	Active    bool
	Items     []items.Item
	Password  string
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

// Create a user in the database
func (user *User) Create() error {
	query := "INSERT INTO users (name, email) VALUES ($1, $2) returning id"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	var password string
	password, err = utilities.HashPassword(user.Password)

	if err != nil {
		log.Fatal(err)
		return err
	}

	err = stmt.QueryRow(user.Name, user.Email, password).Scan(&user.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

//Get is used to fetch a user from the database
func (user *User) Get() error {
	query := "SELECT name, email, ratings FROM users WHERE id = $1 and active = true"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	err = stmt.QueryRow(user.ID).Scan(&user.Name, &user.Email, &user.Ratings, user.CreatedAt)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

//Update a user in the database
func (user *User) Update(changes map[string]string) error {
	query := "UPDATE users SET"
	for k, v := range changes {
		query += " "
		query += k
		query += " "
		query += "="
		query += " "
		query += v
	}

	query += "WHERE id = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(user.ID)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

//Delete a user from the database
func (user *User) Delete() error {
	query := "DELETE FROM users WHERE id = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(user.ID)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

//GetAll users from the database
func (user *User) GetAll() ([]User, error) {
	query := "SELECT name, email, ratings FROM users"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Ratings); err != nil {
			log.Fatal(err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
