package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/usman-174/database"
	"github.com/usman-174/models"
)

var Mykey *models.User

type ErrorResponse struct {
	Error string `json:"error"`
	Msg   string `json:"msg"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := models.User{}
		// Do stuff here
		cookie, err := r.Cookie("token")
		if err != nil {
			fmt.Println("INVALID Token")

			respondWithJSON(w, &ErrorResponse{err.Error(), "Please Login Again"})
			return
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
		if err != nil {
			fmt.Println("INVALID Token")

			respondWithJSON(w, map[string]string{
				"Error": "Please login",
				"Msg":   err.Error(),
			})
			return
		}
		db := database.ConnectDataBase()
		claims := token.Claims.(*jwt.StandardClaims)
		err = db.Preload("Posts").First(&user, "id = ?", claims.Issuer).Error
		if err != nil {

			respondWithJSON(w, map[string]string{
				"Error": "Please login and try again",
				"Msg":   err.Error(),
			})
			return
		}
		if user.ID == 0 {
			respondWithJSON(w, map[string]string{
				"Error": "Please login and try again.",
			})
			return
		}

		ctx = context.WithValue(ctx, Mykey, user)
		// Call the next handler, which can be another middleware in the chain, or the final handler.

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
