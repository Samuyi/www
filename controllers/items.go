package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Samuyi/www/email"
	"github.com/Samuyi/www/models/comments"
	"github.com/Samuyi/www/models/items"
	"github.com/gorilla/websocket"
)

const baseURL = ""

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//CreateItem creates an item
func CreateItem(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("sessionID")
	user, err := getUserFromSession(sessionID)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if !user.Active {
		msg := map[string]string{"error": "Sorry your account isn't activated yet"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if r.Body == nil {
		msg := map[string]string{"error": "Please supply Name, UserID, PhoneNo, LocationID and Instruction"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var item = &items.Item{}
	err = json.NewDecoder(r.Body).Decode(&item)

	log.Println(item.Location)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Please supply a valid name, email and password"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}
	item.UserID = user.ID
	errors := item.Validate()

	if len(errors) > 0 {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)

		return
	}

	err = item.Create()

	if err != nil {
		if err.Error() == "pq: insert or update on table \"items\" violates foreign key constraint \"items_city_fkey\"" {
			msg := map[string]string{"error": "Please supply a valid city"}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(msg)

			return
		}
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	resp := map[string]string{
		"message": "success",
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	return

}

//GetItem gets an item
func GetItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		msg := map[string]string{"error": "id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var item = &items.Item{}
	item.ID = id

	err := item.Get()

	if err != nil {
		msg := map[string]string{"error": "Please try again later"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var comment comments.Comment

	comment.ItemID = id

	itemComments, err := comment.GetItemComments()

	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)

		return
	}

	item.Comments = itemComments

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)

	return
}

//GetItemsInALocation gets all items in a particular location not yet closed
func GetItemsInALocation(w http.ResponseWriter, r *http.Request) {
	locationID := r.URL.Query().Get("location_id")

	if locationID == "" {
		msg := map[string]string{"error": "location_id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Please try again later"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	go func(conn *websocket.Conn) {
		c := time.Tick(2 * time.Minute)

		for range c {

			var itemArray *items.Item
			itemArray.Location.LocationID = locationID

			resp, err := itemArray.ItemsInALocation()
			if err != nil {
				conn.Close()

				return
			}

			err = conn.WriteJSON(resp)

			if err != nil {
				log.Println(err)
				conn.Close()

				return
			}

		}
	}(conn)

}

//GetAllItems gets all items not yet closed
func GetAllItems(w http.ResponseWriter, r *http.Request) {
	var itemsList *items.Item
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Please try again later"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	go func(conn *websocket.Conn) {
		c := time.Tick(2 * time.Minute)

		for range c {
			resp, err := itemsList.GetAllItems()

			if err != nil {
				conn.Close()
			}

			err = conn.WriteJSON(resp)

			if err != nil {
				log.Println(err)
				conn.Close()

				return
			}

		}
	}(conn)

}

//BidItem bids for an item not yet closed
func BidItem(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("sessionID")
	user, err := getUserFromSession(sessionID)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if !user.Active {
		msg := map[string]string{"error": "Sorry your account isn't activated yet"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msg)

		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		msg := map[string]string{"error": "id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if r.Body == nil {
		msg := map[string]string{"error": "Please supply reasons for your bid"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var item = &items.Item{ID: id}

	err = item.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if item.Closed {
		msg := map[string]string{"error": "Sorry item has been closed"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var msg = make(map[string]interface{})

	err = json.NewDecoder(r.Body).Decode(&msg)

	if err != nil {
		log.Println(err)
		res := map[string]string{"error": "Please supply a valid message"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)

		return
	}
	key := id + ":bids"
	_, err = client.HSet(key, user.DisplayName, msg["message"]).Result()

	if err != nil {
		log.Println(err)
		res := map[string]string{"error": "Please try again later there was an error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)

		return
	}

	var mail = &email.Mail{To: item.UserEmail}

	go mail.SendBidAlertMail(item.DisplayName, baseURL+"/?id="+id)

	res := map[string]string{
		"message": "success",
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

	return
}

//GetBidsOnItem gets all bids for an item
func GetBidsOnItem(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("sessionID")
	user, err := getUserFromSession(sessionID)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if !user.Active {
		msg := map[string]string{"error": "Sorry your account isn't activated yet"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msg)

		return
	}

	id := r.URL.Query().Get("id")

	log.Println("i'm here")

	if id == "" {
		msg := map[string]string{"error": "id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var item = &items.Item{ID: id}

	err = item.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if user.ID != item.UserID {
		msg := map[string]string{"error": "Sorry you're not authorized to view this page"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)

		return
	}

	key := id + ":bids"
	resp, err := client.HGetAll(key).Result()

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

	return
}

//CloseItem closes an item from people biding
func CloseItem(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("sessionID")
	user, err := getUserFromSession(sessionID)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if !user.Active {
		msg := map[string]string{"error": "Sorry your account isn't activated yet"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msg)

		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		msg := map[string]string{"error": "id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var item = &items.Item{ID: id}

	err = item.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if user.ID != item.UserID {
		msg := map[string]string{"error": "Sorry you're not authorized to view this page"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)

		return
	}

	item.Closed = true

	err = item.Update()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	return
}

//UpdateItem updates an item
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("sessionID")
	user, err := getUserFromSession(sessionID)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if !user.Active {
		msg := map[string]string{"error": "Sorry your account isn't activated yet"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msg)

		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		msg := map[string]string{"error": "id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var item = &items.Item{ID: id}

	err = item.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if user.ID != item.UserID {
		msg := map[string]string{"error": "Sorry you're not authorized to make a change here"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)

		return
	}
	if item.Closed {
		msg := map[string]string{"error": "Sorry item is closed already"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = item.Update()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	msg := map[string]string{"message": "Success!"}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)

	return
}
