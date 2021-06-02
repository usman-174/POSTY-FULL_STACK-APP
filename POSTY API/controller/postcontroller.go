package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/usman-174/database"
	"github.com/usman-174/models"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var err error
	user := r.Context().Value(Mykey).(models.User)
	post := models.Post{}
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		fmt.Println(err)
		respondWithJSON(w, map[string]string{
			"error": "There was an error,Please try again.",
			"msg":   err.Error(),
		})

		return
	}
	db := database.ConnectDataBase()
	post.UserID = user.ID

	err = db.Create(&post).Error
	if err != nil {
		fmt.Println(err)
		respondWithJSON(w, map[string]string{
			"error": "Couldn't create post.Try again Later.",
			"msg":   err.Error(),
		})

		return
	}
	err = db.Preload("User").Preload("Likes").Find(&post, "id = ?", post.ID).Error
	if err != nil {
		fmt.Println(err)
		respondWithJSON(w, map[string]string{
			"error": "Couldn't create post",
			"msg":   err.Error(),
		})

		return
	}

	respondWithJSON(w, post)
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	var err error
	db := database.ConnectDataBase()
	getPost := []*models.Post{}
	err = db.Preload("User").Preload("Likes").Order("updated_at desc").Find(&getPost).Error
	if err != nil {
		fmt.Println(err)
		respondWithJSON(w, map[string]string{
			"Error": "Could not get posts",
			"Msg":   err.Error(),
		})

		return
	}
	if len(getPost) == 0 {
		respondWithJSON(w, map[string]string{
			"Error": "Not posts found",
		})

		return
	}
	fmt.Println(getPost)
	respondWithJSON(w, &getPost)

}
func GetPost(w http.ResponseWriter, r *http.Request) {
	var err error
	post := models.Post{}

	err = json.NewDecoder(r.Body).Decode(&post)
	db := database.ConnectDataBase()
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"Error": "Invalid Arguments",
			"Msg":   err.Error(),
		})

		return
	}
	err = db.Preload("User").Preload("Likes").Find(&post, "id = ?", post.ID).Error
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"error": "Could not load post. Please refresh",
			"msg":   err.Error(),
		})
		return
	}
	respondWithJSON(w, post)
}
func GetMyPost(w http.ResponseWriter, r *http.Request) {

	var err error
	posts := []models.Post{}
	user := r.Context().Value(Mykey).(models.User)
	db := database.ConnectDataBase()

	err = db.Preload("User").Preload("Likes").Order("updated_at desc").Find(&posts, "user_id = ?", user.ID).Error
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"error": "Could not load post. Please refresh",
			"msg":   err.Error(),
		})
		return
	}
	respondWithJSON(w, posts)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var err error
	user := r.Context().Value(Mykey).(models.User)
	request := map[string]int{}
	err = json.NewDecoder(r.Body).Decode(&request)
	fmt.Printf("The id is = %v\n", request["id"])
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"error": "Invalid arguments",
			"msg":   err.Error(),
		})

		return
	}
	db := database.ConnectDataBase()
	foundPost := &models.Post{}
	err = db.Find(&foundPost, "id = ?", request["id"]).Error
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"error": "Post Not Found.",
			"msg":   err.Error(),
		})

		return
	}
	if foundPost.ID == 0 {
		respondWithJSON(w, map[string]string{
			"error": "Post Not Found.",
		})

		return
	}
	fmt.Println("post.userid = ", foundPost.UserID)
	fmt.Println("user.Id = ", user.ID)
	if foundPost.UserID != user.ID {
		respondWithJSON(w, map[string]string{
			"error": "Only the Post Author can delete the post",
		})

		return
	}
	db.Exec("DELETE FROM posts where id=?", request["id"])

	respondWithJSON(w, map[string]string{
		"msg": "Post Deleted Successfully",
	})

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var err error
	user := r.Context().Value(Mykey).(models.User)
	request := map[string]string{}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"Error": "Invalid arguments 1",
			"Msg":   err.Error(),
		})

		return
	}
	db := database.ConnectDataBase()
	foundPost := &models.Post{}
	postId, err := strconv.Atoi(request["id"])
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"Error": "Cannot convert string id to int id",
			"Msg":   err.Error(),
		})

		return
	}
	err = db.Preload("User").Find(foundPost, "id = ?", postId).Error
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"Error": "Invalid arguments 2",
			"Msg":   err.Error(),
		})

		return
	}

	if foundPost.UserID != user.ID {
		respondWithJSON(w, map[string]string{
			"Error": "Only the Post Author can Update the post",
		})

		return
	}
	foundPost.Body = request["body"]
	foundPost.Title = request["title"]
	err = db.Save(&foundPost).Error
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"Error": "Invalid arguments 3",
			"Msg":   err.Error(),
		})

		return
	}

	respondWithJSON(w, foundPost)

}
