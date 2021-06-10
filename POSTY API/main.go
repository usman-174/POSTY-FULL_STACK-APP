package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/usman-174/app"
)

func main() {
	server := app.Router()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	clientUrl := os.Getenv("CLIENT_URL")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{clientUrl},
		AllowCredentials: true,
	})

	handler := c.Handler(server)
	log.Fatal(http.ListenAndServe(":4000", handler))
	fmt.Println("SERVER UP AND RUNNING")

}
