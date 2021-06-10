package app

import (
	"github.com/gorilla/mux"
	"github.com/usman-174/controller"
	"github.com/usman-174/database"
)

func Router() *mux.Router {
	database.ConnectDataBase()

	router := mux.NewRouter()
	router.HandleFunc("/register", controller.Register).Methods("POST")
	router.HandleFunc("/login", controller.Login).Methods("POST")
	router.HandleFunc("/getposts", controller.GetAllPosts).Methods("GET")
	router.HandleFunc("/getpost", controller.GetPost).Methods("GET")
	protectedRoutes := router.PathPrefix("/x").Subrouter()
	protectedRoutes.Use(controller.AuthMiddleware)
	protectedRoutes.HandleFunc("/user", controller.GetUser).Methods("GET")
	protectedRoutes.HandleFunc("/updateuser", controller.UpdateUser).Methods("POST")
	protectedRoutes.HandleFunc("/post", controller.Post).Methods("POST")
	protectedRoutes.HandleFunc("/myposts", controller.GetMyPost).Methods("GET")
	protectedRoutes.HandleFunc("/logout", controller.Logout).Methods("GET")
	protectedRoutes.HandleFunc("/deletepost", controller.DeletePost).Methods("POST")
	protectedRoutes.HandleFunc("/updatepost", controller.UpdatePost).Methods("POST")
	protectedRoutes.HandleFunc("/likepost", controller.Likepost).Methods("POST")
	return router
}
