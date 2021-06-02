package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/usman-174/app"
)

func main() {
	server := app.Router()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(server)
	log.Fatal(http.ListenAndServe(":4000", handler))
	fmt.Println("SERVER UP AND RUNNING")

}
