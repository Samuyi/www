package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Samuyi/www/utilities"
	jwt "github.com/dgrijalva/jwt-go"
)

//Middleware is a function that runs before a controller
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(methods ...string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			for _, m := range methods {
				if r.Method == m {
					f(w, r)

					return
				}
			}

			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}

var secret = []byte(os.Getenv("AUTH_KEY"))

//Auth is a middleware to authenticate request on the server
func Auth() Middleware {

	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")

			if bearerToken == "" {
				msg := map[string]string{"message": "Please supply a token"}
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(msg)

				return
			}

			tokenString := bearerToken[7:]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return secret, nil
			})

			if err != nil {
				msg := map[string]string{"message": "Sorry token is invalid"}
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(msg)

				return
			}
			if !token.Valid {
				msg := map[string]string{"message": "Sorry token is invalid"}
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(msg)

				return
			}

			errors := token.Claims.Valid()

			if errors != nil {
				msg := map[string]string{"message": "Sorry token is invalid"}
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(msg)

				return
			}

			claims := token.Claims.(jwt.MapClaims)

			sessionID := claims["sessionID"].(string)

			_, err = utilities.GetSession(sessionID)

			if err != nil {
				msg := map[string]string{"message": "Sorry session has expired"}
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(msg)

				return
			}

			r.Header.Set("sessionID", sessionID)
			f(w, r)

		}
	}

}

//WithCors implemets cross origin request middleware
func WithCors() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			fn(w, r)

			return
		}
	}

}

//ChainMiddlewares chains one or two middleware together
func ChainMiddlewares(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {

	for _, m := range middlewares {
		f = m(f)
	}
	return f

}
