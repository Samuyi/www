package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Samuyi/www/models/locations"
)

//CreateLocation creates a location
func CreateLocation(w http.ResponseWriter, r *http.Request) {
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

	var location = &locations.Location{}

	if r.Body == nil {
		msg := map[string]string{"error": "Please supply a name, state and country code of a location"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
	}

	err = json.NewDecoder(r.Body).Decode(&location)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Please supply a valid name, email and password"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}
	location.UserID = user.ID

	errors := location.Validate()

	if len(errors) > 0 {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)

		return
	}

	err = location.Create()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	msg := map[string]string{"message": "Success!"}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)

	return

}

//GetLocations gets all locations
func GetLocations(w http.ResponseWriter, r *http.Request) {
	var location = &locations.Location{}

	locations, err := location.GetAll()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(locations)

	return
}

//GetLocation gets a location
func GetLocation(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		msg := map[string]string{"error": "id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var location = &locations.Location{}

	location.LocationID = id

	err := location.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	location.UserID = ""

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(location)

}

//UpdateLocation updates a location
func UpdateLocation(w http.ResponseWriter, r *http.Request) {
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
		msg := map[string]string{"error": "Please supply a valid id"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var location = &locations.Location{LocationID: id}

	if r.Body == nil {
		msg := map[string]string{"error": "Please supply values to be updated"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var changes = make(map[string]string)

	err = json.NewDecoder(r.Body).Decode(&changes)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry only city, country or state names can be updated"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = location.Update(changes)

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
