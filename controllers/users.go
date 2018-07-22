package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Samuyi/www/email"
	"github.com/Samuyi/www/models/users"
	"github.com/Samuyi/www/utilities"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
)

//MyCustomClaims is a type for the jwt claims
type MyCustomClaims struct {
	sessionID string
	jwt.StandardClaims
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SIGNING_KEY")))

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func createToken() (map[string]string, error) {
	mySigningKey := []byte(os.Getenv("SIGNING_KEY"))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	letters := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	r.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})

	sessionID := strings.Join(letters[2:12], "")
	exp := time.Now().Unix() + (86400 * 3)

	claims := MyCustomClaims{
		sessionID,
		jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	session := map[string]string{
		"sessionID": sessionID,
		"token":     ss,
	}

	return session, nil
}

func setSession(sessionID string, sessionValues map[interface{}]interface{}) error {

	session := map[string]interface{}{}

	for k, v := range sessionValues {
		key := k.(string)
		value := v.(string)
		session[key] = value
	}

	_, err := client.HMSet(sessionID, session).Result()
	if err != nil {
		log.Println(err)
		return err
	}

	duration, _ := time.ParseDuration("36h") // expire after three days
	_, err = client.Expire(sessionID, duration).Result()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//GetSessiom gets the values of a session
func getSession(sessionID string) (users.User, error) {
	var user users.User

	session, err := client.HGetAll(sessionID).Result()

	if err != nil {
		log.Println(err)
		return user, err
	}

	if session == nil {
		return user, fmt.Errorf("Session has expired")
	}

	user.Active, _ = strconv.ParseBool(session["Active"])
	user.ID = session["userID"]
	user.FirstName = session["FirstName"]
	user.DisplayName = session["DisplayName"]
	user.LastName = session["LastName"]
	user.Avatar = session["Avatar"]

	return user, nil
}

//RegisterUser controller to create a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user users.User

	if r.Body == nil {
		msg := map[string]string{"error": "Please supply name, email and password"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Please supply a valid name, email and password"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	errors := user.Validate()

	if len(errors) > 0 {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)

		return
	}

	err = user.Create()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	user.Password = ""
	id, err := utilities.SetUserForConfirmation(user.ID)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var mail = &email.Mail{To: user.Email}

	err = mail.SendConfirmationMail(user.DisplayName, baseURL+"/?key="+id)
	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	sessionInfo, err := createToken()
	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}
	sessionID := sessionInfo["sessionID"]
	token := sessionInfo["token"]

	session, err := store.Get(r, sessionID)
	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	session.Values["FirstName"] = user.FirstName
	session.Values["LastName"] = user.LastName
	session.Values["DisplayName"] = user.DisplayName
	session.Values["Active"] = user.Active
	session.Values["Avartar"] = user.Avatar
	session.Values["userID"] = user.ID

	err = session.Save(r, w)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = setSession(sessionID, session.Values)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	resp := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	return
}

//Login logs a user into the application
func Login(w http.ResponseWriter, r *http.Request) {
	var user users.User

	if r.Body == nil {
		msg := map[string]string{"error": "Please supply email and password"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Please supply a valid email and password"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	password := user.Password
	err = user.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	ok := utilities.CheckPassword(password, user.Password)

	if !ok {
		msg := map[string]string{"error": "Invalid Password"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	user.Password = ""

	sessionInfo, err := createToken()

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	sessionID := sessionInfo["sessionID"]

	session, err := store.Get(r, sessionID)
	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	session.Values["FirstName"] = user.FirstName
	session.Values["LastName"] = user.LastName
	session.Values["DisplayName"] = user.DisplayName
	session.Values["Active"] = user.Active
	session.Values["Avartar"] = user.Avatar
	session.Values["userID"] = user.ID

	err = session.Save(r, w)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = setSession(sessionID, session.Values)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error please try again later"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	resp := map[string]string{
		"token": sessionInfo["token"],
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

//ConfirmUser confirms a user's email address
func ConfirmUser(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		msg := map[string]string{"error": "Key required"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	id, err := utilities.GetUserForConfirmation(key)

	if err != nil {
		msg := map[string]interface{}{"error": err}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var user = &users.User{ID: id}

	err = user.Get()

	if err != nil {
		msg := map[string]string{"error": "Please try again later"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = user.SetUserActive()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error. Please try again latter"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	sessionInfo, err := createToken()
	sessionID := sessionInfo["sessionID"]
	token := sessionInfo["token"]
	session, err := store.Get(r, sessionID)

	session.Values["FirstName"] = user.FirstName
	session.Values["LastName"] = user.LastName
	session.Values["DisplayName"] = user.DisplayName
	session.Values["Active"] = user.Active
	session.Values["Avartar"] = user.Avatar
	session.Values["userID"] = user.ID

	err = session.Save(r, w)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error. Please try again latter"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = setSession(sessionID, session.Values)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	resp := map[string]string{
		"token": token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

	return

}

//ForgotPassword changes a user's password if they forgot their password
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		msg := map[string]string{"error": "Please supply a valid email"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var user *users.User

	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}
	user.ID = ""
	err = user.GetID()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if user.ID == "" {
		msg := map[string]string{"error": "Sorry this email doesn't exist."}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	letters := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	random.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})

	password := strings.Join(letters[2:12], "")

	user.Password = password
	err = user.UpdatePassword()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	// send email to user of updated paasword in the database and ask user to update their password

	var mail = &email.Mail{To: user.Email}

	err = mail.EmailPassword(password, baseURL+"/login")

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	resp := map[string]string{
		"message": "We have sent a tempoary password to your email. Please check your email.",
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

//UpdateUser updates a user's data in the database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("sessionID")

	user, err := getSession(sessionID)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	if r.Body == nil {
		msg := map[string]string{"error": "Please supply values to be updated"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)

		return
	}

	user.Password = ""

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Println(err)
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = user.Update()

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

//DeleteUser Deletes a user from the application
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("sessionID")

	user, err := getSession(sessionID)

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	err = user.Delete()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	_, err = client.Del(sessionID).Result()

	if err != nil {
		log.Println(err)
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

//GetAllUsers fetches all users from the application
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var user = &users.User{}

	users, err := user.GetAll()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

	return

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id != "" {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	var user = &users.User{ID: id}

	err := user.Get()

	if err != nil {
		msg := map[string]string{"error": "Sorry there was an internal server error"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

	return

}

//LogOut logs a user out of the application
func LogOut(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("sessionID")

	_, err := client.Del(sessionID).Result()

	if err != nil {
		log.Println(err)
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
