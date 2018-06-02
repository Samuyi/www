package items

import (
	"database/sql"
	"log"
	"time"

	"github.com/Samuyi/www/models/comments"

	_ "github.com/lib/pq" // postgres driver
)

//Item data structure
type Item struct {
	ID        string
	Name      string
	UserID    string
	Location  string
	Closed    bool
	Comments  []comments.Comment
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

//Create an item in the databsae
func (item *Item) Create() error {
	query := "INSERT INTO items (user_id, location, ) VALUES ($1, $2, $3) returning id"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	err = stmt.QueryRow(item.Name, item.Location, item.UserID).Scan(&item.ID)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

//Get an item from the database
func (item *Item) Get() error {
	query := "SELECT name, user_id, location, created_at FROM items  where id = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	err = stmt.QueryRow(item.ID).Scan(&item.Name, item.UserID, &item.Location, &item.CreatedAt)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return err
}

//GetAll items belonging to a particular user
func (item *Item) GetAll(userID string) ([]Item, error) {
	query := "SELECT name, user_id, location, created_at FROM items where user_id = $1"

	stmt, err := db.Prepare(query)
	stmt.Close()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	rows, err := stmt.Query(userID)

	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var items []Item

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Location, &item.Closed, &item.CreatedAt); err != nil {
			log.Fatal(err)
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return items, err

}

//Update an item in the database
func (item *Item) Update() error {
	query := "UPDATE items SET name = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(item.ID)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return err
}

//Delete itemm from the database
func (item *Item) Delete() error {
	query := "DELETE FROM items where id = $id"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(item.ID)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return err
}
