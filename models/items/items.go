package items

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Samuyi/www/models/comments"
	"github.com/Samuyi/www/models/locations"
	validator "github.com/asaskevich/govalidator"
	_ "github.com/lib/pq" // postgres driver
)

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

//Item data structure
type Item struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	UserID      string             `json:"user_id"`
	DisplayName string             `json:"display_name"`
	UserEmail   string             `json:"user_email"`
	PhoneNo     string             `json:"phone_no"`
	Location    locations.Location `json:"location"`
	Closed      bool               `json:"closed"`
	Instruction string             `json:"instruction"`
	Comments    []comments.Comment `json:"comments,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty"`
}

//Validate item struct
func (item *Item) Validate() map[string]string {
	var errors = make(map[string]string)

	if len(item.Name) <= 2 {
		message := "item name must be at least three words"
		errors["Invalid item Name"] = message
	}

	if !validator.IsUUID(item.UserID) {
		message := "Please supply a valid user id"
		errors["Invalid UserID"] = message
	}

	if len(item.PhoneNo) <= 5 {
		message := "Please supply a valid phone number"
		errors["Invalid phone number"] = message
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

//Create an item in the databsae
func (item *Item) Create() error {
	query := "INSERT INTO items (user_id, name, phone_no, instruction, city ) VALUES ($1, $2, $3, $4, $5) returning id"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(item.UserID, item.Name, item.PhoneNo, item.Instruction, item.Location.City).Scan(&item.ID)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Get an item from the database
func (item *Item) Get() error {
	query := "SELECT name, display_name, email, items.user_id, items.city, instruction, phone_no, items.created_at as created_at, state, country FROM items INNER JOIN users ON items.user_id = users.id INNER JOIN locations ON locations.city = items.city where items.id = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(item.ID).Scan(&item.Name, &item.DisplayName, &item.UserEmail, &item.UserID, &item.Location.City, &item.Instruction, &item.PhoneNo, &item.CreatedAt, &item.Location.State, &item.Location.Country)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Update an item in the database
func (item *Item) Update() error {
	query := "UPDATE items SET name = $1, phone_no = $2, closed = $3, instruction = $4, updated_at=$5 where id = $6"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	item.UpdatedAt = time.Now()
	_, err = stmt.Exec(item.Name, item.PhoneNo, item.Closed, item.Instruction, item.UpdatedAt, item.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Delete itemm from the database
func (item *Item) Delete() error {
	query := "DELETE FROM items where id = $id"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(item.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//ItemsInALocation gets items in a particular location
func (item *Item) ItemsInALocation() ([]Item, error) {
	query := "SELECT id, name, user_id, user.display_name, email, phone_no, instruction FROM items INNER JOIN users ON items.user_id = users.id WHERE location_id = $1 and closed = false ORDER BY created_at DESC"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query(item.Location.LocationID)

	var itemArray []Item

	defer rows.Close()
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.UserID, &item.DisplayName, &item.UserEmail, &item.PhoneNo, &item.Instruction, &item.CreatedAt); err != nil {
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

//GetAllItems gets all items still open currently
func (item *Item) GetAllItems() ([]Item, error) {
	query := "SELECT id, name, user_id, user.display_name, email, phone_no, instruction, city, created_at FROM items INNER JOIN users ON items.user_id = user.id WHERE closed = false ORDER BY created_at DESC"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()
	var itemArray []Item

	defer rows.Close()

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.UserID, &item.PhoneNo, &item.Instruction, &item.Location.City, &item.CreatedAt); err != nil {
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
