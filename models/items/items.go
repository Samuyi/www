package items

import (
	"database/sql"
	"log"
	"time"

	"github.com/Samuyi/www/models/comments"
	"github.com/Samuyi/www/models/locations"
	validator "github.com/asaskevich/govalidator"
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

//Item data structure
type Item struct {
	ID          string
	Name        string
	UserID      string
	UserName    string
	UserEmail   string
	PhoneNo     int
	Location    locations.Location
	Closed      bool
	Instruction string
	Comments    []comments.Comment
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Validate item struct
func (item *Item) Validate() map[string]string {
	var errors map[string]string

	if len(item.Name) <= 2 {
		message := "item name must be at least three words"
		errors["Invalid item Name"] = message
	}

	if !validator.IsUUID(item.UserID) {
		message := "Please supply a valid user id"
		errors["Invalid UserID"] = message
	}

	if !validator.IsUUID(item.Location.LocationID) {
		message := "item must have a valid location ID"
		errors["Invalid LocationID"] = message
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

//Create an item in the databsae
func (item *Item) Create() error {
	query := "INSERT INTO items (user_id, name, phone_no, instruction, location_id, ) VALUES ($1, $2, $3, $4, $5) returning id"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(item.UserID, item.Name, item.PhoneNo, item.Instruction, item.Location.LocationID).Scan(&item.ID)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Get an item from the database
func (item *Item) Get() error {
	query := "SELECT name, user.first_name, email, user_id, location_id, instruction, created_at FROM items INNER JOIN ON items.user_id = users.id  where id = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(item.ID).Scan(&item.Name, &item.UserName, &item.UserEmail, item.UserID, &item.Location.LocationID, &item.Instruction, &item.CreatedAt)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Update an item in the database
func (item *Item) Update() error {
	query := "UPDATE items SET name = $1, phone_no = $2, closed = $3, instruction = $4 where id = $5"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(item.Name, item.PhoneNo, item.Closed, item.Instruction, item.ID)

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

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.UserID, &item.UserName, &item.UserEmail, &item.PhoneNo, &item.Instruction, &item.CreatedAt); err != nil {
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
	query := "SELECT id, name, user_id, user.display_name, email, phone_no, instruction FROM items INNER JOIN users ON items.user_id = user.id WHERE closed = false ORDER BY created_at DESC"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()
	var itemArray []Item

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.UserID, &item.PhoneNo, &item.Instruction, &item.CreatedAt); err != nil {
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
