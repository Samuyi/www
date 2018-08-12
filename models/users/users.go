package users

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Samuyi/www/models/items"
	utilities "github.com/Samuyi/www/utilities"
	validate "github.com/asaskevich/govalidator"
	_ "github.com/lib/pq" // postgres driver
)

//db is the database connector
var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "help"
	password = "help"
	dbname   = "help.ng"
)

func init() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println(err)
	}
	err = db.Ping()

	if err != nil {
		log.Println(err)
	}

	log.Println("connected to database")
}

//User data structure
type User struct {
	ID          string       `json:"id"`
	DisplayName string       `json:"display_name"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	Email       string       `json:"email"`
	Ratings     int          `json:"ratings,omitempty"`
	Avatar      string       `json:"avatar,omitempty"`
	Active      bool         `json:"active"`
	Items       []items.Item `json:"items,omitempty"`
	Password    string       `json:"password,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty"`
}

//Validate the fields of a user
func (user *User) Validate() map[string]string {
	var errors = make(map[string]string)

	if len(user.Password) < 8 {
		message := "Password must be greater than 7 characters."
		errors["Password Error"] = message
	}

	if !validate.IsEmail(user.Email) {
		message := "Please supply a valid email"
		errors["Email Error"] = message
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Create a user in the database
func (user *User) Create() error {
	query := "INSERT INTO users (first_name, last_name, display_name, email, password, avatar) VALUES ($1, $2, $3, $4, $5, $6) returning id;"
	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	var password string
	password, err = utilities.HashPassword(user.Password)

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(user.FirstName, user.LastName, user.DisplayName, user.Email, password, user.Avatar).Scan(&user.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Get is used to fetch a user from the database
func (user *User) Get() error {
	query := "SELECT first_name, last_name, display_name, email, ratings, active, password, created_at FROM users WHERE id = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(user.ID).Scan(&user.FirstName, &user.LastName, &user.DisplayName, &user.Email, &user.Ratings, &user.Active, &user.Password, &user.CreatedAt)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//GetUserByName gets a users based on username
func (user *User) GetUserByName() error {
	query := "SELECT first_name, last_name, email, ratings, active  created_at FROM users WHERE display_name = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(user.DisplayName).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Ratings, &user.Active, &user.CreatedAt)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(user)
	return nil
}

//GetID gets the password asociated with an email
func (user *User) GetID() error {
	query := "SELECT id, password, active, display_name, first_name, last_name, avatar FROM users where email = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(user.Email).Scan(&user.ID, &user.Password, &user.Active, &user.DisplayName, &user.FirstName, &user.LastName, &user.Avatar)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Update a user in the database
func (user *User) Update() error {
	user.UpdatedAt = time.Now()
	if user.Password == "" {
		query := "UPDATE users SET first_name = $1, last_name = $2, display_name = $3, updated_at=$4 WHERE id = $5"

		stmt, err := db.Prepare(query)
		defer stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}

		_, err = stmt.Exec(user.FirstName, user.LastName, user.DisplayName, user.UpdatedAt, user.ID)

		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
	query := "UPDATE users SET first_name = $1, last_name = $2, display_name = $3, password = $4, updated_at=$5 WHERE id = $6"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	var password string
	password, err = utilities.HashPassword(user.Password)

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(user.FirstName, user.LastName, user.DisplayName, password, user.UpdatedAt, user.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

//UpdatePassword of a user
func (user *User) UpdatePassword() error {
	query := "UPDATE users SET password = $1, updated_at=$2 where id = $3"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	var password string
	password, err = utilities.HashPassword(user.Password)

	if err != nil {
		log.Println(err)
		return err
	}
	user.UpdatedAt = time.Now()

	_, err = stmt.Exec(password, user.UpdatedAt, user.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//SetUserActive makes a user active on the network
func (user *User) SetUserActive() error {
	query := "UPDATE users SET active = true, updated_at=$1 WHERE id = $2"

	stmt, err := db.Prepare(query)

	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	user.UpdatedAt = time.Now()
	_, err = stmt.Exec(user.UpdatedAt, user.ID)

	user.Active = true

	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(user.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//GetAll users from the database
func GetAll() ([]User, error) {
	query := "SELECT id, display_name, email, ratings, avatar FROM users ORDER BY created_at DESC"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var users []User

	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.DisplayName, &user.Email, &user.Ratings, &user.Avatar); err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

//GetAllItems gets all items belonging to a userbelonging to a particular user
func (user *User) GetAllItems() ([]items.Item, error) {
	query := "SELECT id, name, location_id, instruction, closed, created_at FROM items where user_id = $1 and active = true"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query(user.ID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var itemArray []items.Item

	defer rows.Close()

	for rows.Next() {
		var item items.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Location.LocationID, &item.Instruction, &item.Closed, &item.CreatedAt); err != nil {
			log.Println(err)
			return nil, err
		}
		itemArray = append(itemArray, item)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return itemArray, nil
}
