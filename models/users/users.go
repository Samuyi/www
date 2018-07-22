package users

import (
	"database/sql"
	"log"
	"time"

	"github.com/Samuyi/www/models/items"
	utilities "github.com/Samuyi/www/utilities"
	validate "github.com/asaskevich/govalidator"
	_ "github.com/lib/pq" // postgres driver
)

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("postgres", "user=help.ng passowrd=this is the password for help.ng sslmode=disable")
	if err != nil {
		log.Println(err)
	}

}

//User data structure
type User struct {
	ID          string
	DisplayName string
	FirstName   string
	LastName    string
	Email       string
	Ratings     string
	Avatar      string
	Active      bool
	Items       []items.Item
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Validate the fields of a user
func (user *User) Validate() map[string]string {
	var errors map[string]string

	if len(user.Password) < 8 {
		message := "Password must be greater than 7 characters."
		errors["Password Error"] = message
	}

	if validate.IsEmail(user.Email) {
		message := "Email is not valid"
		errors["Email Error"] = message
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Create a user in the database
func (user *User) Create() error {
	query := "INSERT INTO users (first_name, last_name, display_name, email, password, avatar) VALUES ($1, $2, $3, $4, $5, $6) returning id"

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
	query := "SELECT fisrt_name, last_name display_name email, ratings, password FROM users WHERE id = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(user.ID).Scan(&user.FirstName, &user.LastName, &user.DisplayName, &user.Email, &user.Ratings, &user.Password)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//GetID gets the password asociated with an email
func (user *User) GetID() error {
	query := "SELECT id FROM users where email = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(user.Email).Scan(&user.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Update a user in the database
func (user *User) Update() error {
	if user.Password == "" {
		query := "UPDATE users SET first_name = $1, last_name = $2, display_name = $3 WHERE id = $4"

		stmt, err := db.Prepare(query)
		defer stmt.Close()

		if err != nil {
			log.Println(err)
			return err
		}

		_, err = stmt.Exec(user.FirstName, user.LastName, user.DisplayName, user.ID)

		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
	query := "UPDATE users SET first_name = $1, last_name = $2, display_name = $3, password = $4 WHERE id = $5"
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

	_, err = stmt.Exec(user.FirstName, user.LastName, user.DisplayName, password, user.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

//UpdatePassword of a user
func (user *User) UpdatePassword() error {
	query := "UPDATE users SET password = $1 where id = $2"

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

	_, err = stmt.Exec(password, user.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//SetUserActive makes a user active on the network
func (user *User) SetUserActive() error {
	query := "UPDATE user SET active = true WHERE id = $1"

	stmt, err := db.Prepare(query)

	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(user.ID)

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
func (user *User) GetAll() ([]User, error) {
	query := "SELECT id, display_name, email, ratings, avatar FROM users WHERE active = true"

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

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.DisplayName, &user.Email, &user.Ratings, user.Avatar); err != nil {
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
	stmt.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query(user.ID)

	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var itemArray []items.Item

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
