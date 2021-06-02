package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/usman-174/database"
	"github.com/usman-174/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginData struct {
	Email    string
	Password string
}

const SecretKey = "secret"

var Users []models.User

func Register(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDataBase()
	user := models.User{}
	reqBody := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		fmt.Println(err)
		respondWithJSON(w, map[string]string{
			"error": "Registration Failed",
			"msg":   err.Error(),
		})

		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(reqBody["password"]), 8)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"msg": err.Error(),
		})

		return
	}

	db.Find(&user, "email = ?", reqBody["email"])
	if user.ID != 0 {
		respondWithJSON(w, map[string]string{
			"error": "User with this email already exists",
		})
		return
	}
	user.Password = string(password)
	user.Email = reqBody["email"]
	user.Username = reqBody["username"]
	err = db.Create(&user).Error
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": err.Error(),
		})

		return
	}
	respondWithJSON(w, user)
}

func empty(s string) bool {
	return len(strings.TrimSpace(s)) < 4
}

func Login(w http.ResponseWriter, r *http.Request) {

	var err error
	reqBody := LoginData{}
	user := models.User{}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		fmt.Println("error =1", err)
		respondWithJSON(w, map[string]string{
			"error": "Login Failed",
			"msg":   err.Error(),
		})
		return
	}

	emptyemail := empty(reqBody.Email)
	emptypass := empty(reqBody.Password)
	if emptyemail {
		respondWithJSON(w, map[string]string{
			"error": "Please enter valid email",
		})
		return
	}
	if emptypass {
		respondWithJSON(w, map[string]string{
			"error": "Please enter a valid password",
		})
		return
	}
	db := database.ConnectDataBase()
	err = db.First(&user, "email = ?", reqBody.Email).Error
	if err != nil {
		fmt.Println("error =2", err)
		respondWithJSON(w, map[string]string{
			"error": "User not found.",
			"msg":   err.Error(),
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		fmt.Println("error =3", err)
		respondWithJSON(w, map[string]string{
			"error": "Invalid Password",
			"msg":   err.Error(),
		})
		return
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		fmt.Println("error =4", err)
		respondWithJSON(w, map[string]string{
			"error": "Could not login",
			"msg":   err.Error(),
		})

		return
	}
	cookie := http.Cookie{Name: "token", Value: token, HttpOnly: true, Path: "/", Expires: time.Now().Add(time.Hour * 24)}
	http.SetCookie(w, &cookie)
	respondWithJSON(w, map[string]string{
		"msg": "You are logged In",
	})
}
func GetUser(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(Mykey).(models.User)
	respondWithJSON(w, user)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOGGIG YOU OUT")
	cookie := http.Cookie{Name: "token", Value: "", HttpOnly: true, Path: "/", Expires: time.Now().Add(-time.Hour)}

	http.SetCookie(w, &cookie)
	respondWithJSON(w, map[string]string{
		"msg": "You are logged out now.",
	})
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(Mykey).(models.User)
	requestBody := map[string]string{}
	err1 := json.NewDecoder(r.Body).Decode(&requestBody)
	if err1 != nil {
		fmt.Println("error =1", err1)
		respondWithJSON(w, map[string]string{
			"error": "There was an error updating Profile.Try again",
			"msg":   err1.Error(),
		})
		return
	}
	fmt.Println(requestBody["email"], requestBody["username"], requestBody["newpassword"],
		requestBody["oldpassword"])
	reqUserId, erra := strconv.Atoi(requestBody["id"])
	if erra != nil {
		respondWithJSON(w, map[string]string{
			"error": "There was an error please try again",
		})
		return

	}
	if user.ID != uint(reqUserId) {
		respondWithJSON(w, map[string]string{
			"error": "There was an error please Login again and then try",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody["oldpassword"]))
	if err != nil {
		fmt.Println("error =3", err)
		respondWithJSON(w, map[string]string{
			"error": "Invalid Old Password",
			"msg":   err.Error(),
		})
		return
	}
	if requestBody["newpassword"] != "" {
		if requestBody["newpassword"] == requestBody["oldpassword"] {
			respondWithJSON(w, map[string]string{
				"error": "New Password cannot be same as it was before.",
			})

			return
		}
		password, err := bcrypt.GenerateFromPassword([]byte(requestBody["newpassword"]), 8)

		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"msg": err.Error(),
			})

			return
		}
		user.Password = string(password)
	}
	if user.Email == requestBody["email"] && user.Username == requestBody["username"] && requestBody["newpassword"] == requestBody["oldpassword"] {
		respondWithJSON(w, map[string]string{
			"error": "Please update atleast one of the fields",
		})

		return
	} else if len(requestBody["username"]) < 3 {
		respondWithJSON(w, map[string]string{
			"error": "Username length must be atleast 3 characters long",
		})
		return

	} else if !strings.Contains(requestBody["email"], "@") {
		respondWithJSON(w, map[string]string{
			"error": "Enter a valid Email",
		})
		return

	}
	db := database.ConnectDataBase()
	user.Email = requestBody["email"]
	user.Username = requestBody["username"]
	errx := db.Preload("Posts").Save(&user).Error
	if errx != nil {
		fmt.Println("db save error=", errx.Error())
		respondWithJSON(w, map[string]string{
			"error": "Enter a valid Email",
			"msg":   errx.Error(),
		})
		return
	}
	respondWithJSON(w, user)

}
func respondWithJSON(w http.ResponseWriter, payload interface{}) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
