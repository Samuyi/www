package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Samuyi/www/models/comments"
)

//CreateComment creates a comment
func CreateComment(w http.ResponseWriter, r *http.Request) {
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
		msg := map[string]string{"error": "Sorry you need to supply an item id and a comment text"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var comment comments.Comment

	err = json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Please supply a valid email and password"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	comment.Username = user.DisplayName

	err = comment.Create()

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

//CreateReply creates a reply
func CreateReply(w http.ResponseWriter, r *http.Request) {
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

	params := mux.Vars(r)
	commentID := params["comment_id"]

	if commentID == "" {
		msg := map[string]string{"error": "comment id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if r.Body == nil {
		msg := map[string]string{"error": "Sorry you need to supply an item id and a comment text"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var reply comments.Reply

	err = json.NewDecoder(r.Body).Decode(&reply)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	reply.CommentID = commentID
	reply.Username = user.DisplayName

	err = reply.Create()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an error"}
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

//GetComment gets a comment
func GetComment(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		msg := map[string]string{"error": "id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var comment comments.Comment

	comment.ID = id

	err := comment.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = comment.GetReplies()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)

	return

}

//GetReplies gets all replies for a comment
func GetReplies(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	commentID := params["comment_id"]

	if commentID == "" {
		msg := map[string]string{"error": "comment id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}
	var comment comments.Comment

	comment.ID = commentID

	err := comment.GetReplies()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment.Replies)

	return

}

//GetItemComments gets all comments for an item
func GetItemComments(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		msg := map[string]string{"error": "id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var comment comments.Comment

	comment.ItemID = id

	commentArray, err := comment.GetItemComments()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(commentArray)

	return

}

//UpdateComment updates a comment
func UpdateComment(w http.ResponseWriter, r *http.Request) {
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

	var comment = &comments.Comment{ID: id}

	err = comment.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if user.DisplayName != comment.Username {
		msg := map[string]string{"error": "Sorry you're not authorized to view this page"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	comment.Replies = []comments.Reply{}

	err = comment.Update()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)

}

//UpdateReply updates a reply
func UpdateReply(w http.ResponseWriter, r *http.Request) {
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

	var reply = &comments.Reply{ID: id}

	err = reply.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if user.DisplayName != reply.Username {
		msg := map[string]string{"error": "Sorry you're not authorized to view this page"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = json.NewDecoder(r.Body).Decode(&reply)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = reply.Update()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reply)

}

//DeleteComment deletes a comment
func DeleteComment(w http.ResponseWriter, r *http.Request) {
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

	var comment = &comments.Comment{ID: id}

	err = comment.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if user.DisplayName != comment.Username {
		msg := map[string]string{"error": "Sorry you're not authorized to carry out this activity"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = comment.Delete()

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

//DeleteReply deletes a reply
func DeleteReply(w http.ResponseWriter, r *http.Request) {
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

	params := mux.Vars(r)
	commentID := params["comment_id"]

	if commentID == "" {
		msg := map[string]string{"error": "comment id required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var reply = &comments.Reply{ID: id, CommentID: commentID}

	err = reply.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if user.DisplayName != reply.Username {
		msg := map[string]string{"error": "Sorry you're not authorized to view this page"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = reply.Delete()

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
