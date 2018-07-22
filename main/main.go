package main

import (
	"net/http"

	"github.com/Samuyi/www/controllers"
	"github.com/Samuyi/www/middleware"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/users", middleware.ChainMiddlewares(controllers.RegisterUser, middleware.Method("POST, OPTIONS"), middleware.WithCors())).Methods("POST, OPTIONS")
	router.HandleFunc("/api/users", middleware.ChainMiddlewares(controllers.GetAllUsers, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/users/user", middleware.ChainMiddlewares(controllers.GetUser, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/login", middleware.ChainMiddlewares(controllers.Login, middleware.Method("POST", "OPTIONS"))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/logout", middleware.ChainMiddlewares(controllers.LogOut, middleware.Method("GET"), middleware.Auth())).Methods("GET")
	router.HandleFunc("/api/confirm-email", middleware.ChainMiddlewares(controllers.ConfirmUser, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/users", middleware.ChainMiddlewares(controllers.UpdateUser, middleware.Method("PUT", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/users", middleware.ChainMiddlewares(controllers.DeleteUser, middleware.Method("DELETE", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/forgot-password", middleware.ChainMiddlewares(controllers.ForgotPassword, middleware.Method("POST", "OPTIONS"), middleware.WithCors())).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/items", middleware.ChainMiddlewares(controllers.CreateItem, middleware.Method("POST", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/items/{id}", middleware.ChainMiddlewares(controllers.CreateItem, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/items", middleware.ChainMiddlewares(controllers.GetAllItems, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/items/location", middleware.ChainMiddlewares(controllers.GetItemsInALocation, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/items", middleware.ChainMiddlewares(controllers.UpdateItem, middleware.Method("PUT", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/items", middleware.ChainMiddlewares(controllers.CloseItem, middleware.Method("PATCH", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/api/items/bid", middleware.ChainMiddlewares(controllers.BidItem, middleware.Method("POST", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/items/bid", middleware.ChainMiddlewares(controllers.GetBidsOnItem, middleware.Method("GET"), middleware.Auth())).Methods("GET")

	router.HandleFunc("/api/comments", middleware.ChainMiddlewares(controllers.CreateComment, middleware.Method("POST", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/comments", middleware.ChainMiddlewares(controllers.GetComment, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/comments/item", middleware.ChainMiddlewares(controllers.GetItemComments, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/comments", middleware.ChainMiddlewares(controllers.UpdateComment, middleware.Method("PUT", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/comments", middleware.ChainMiddlewares(controllers.DeleteComment, middleware.Method("DELETE", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/comments/{comment_id}/reply", middleware.ChainMiddlewares(controllers.CreateReply, middleware.Method("POST", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/comments/{comment_id}/reply", middleware.ChainMiddlewares(controllers.UpdateReply, middleware.Method("PUT", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/comments/{comment_id}/reply", middleware.ChainMiddlewares(controllers.DeleteReply, middleware.Method("DELETE", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/locations", middleware.ChainMiddlewares(controllers.CreateLocation, middleware.Method("POST", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/locations", middleware.ChainMiddlewares(controllers.GetLocations, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/locations/location/", middleware.ChainMiddlewares(controllers.CreateLocation, middleware.Method("GET"))).Methods("GET")
	router.HandleFunc("/api/locations", middleware.ChainMiddlewares(controllers.CreateLocation, middleware.Method("DELETE", "OPTIONS"), middleware.WithCors(), middleware.Auth())).Methods("DELETE", "OPTIONS")

	http.Handle("/", router)

	http.ListenAndServe("8080", nil)

}
