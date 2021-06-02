package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/usman-174/database"
	"github.com/usman-174/models"
)

func Likepost(w http.ResponseWriter, r *http.Request) {
	var err error
	user := r.Context().Value(Mykey).(models.User)
	request := map[string]uint{}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"Error": "Invalid arguments",
			"Msg":   err.Error(),
		})

		return
	}
	db := database.ConnectDataBase()

	// postId, err := strconv.Atoi(request["id"])
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	respondWithJSON(w, map[string]string{
	// 		"Error": "Invalid req.body",
	// 		"Msg":   err.Error(),
	// 	})
	//
	// 	return
	// }
	foundpost := models.Post{}
	err = db.Find(&foundpost, "id = ?", request["id"]).Error
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"error": "Post not found",
			"msg":   err.Error(),
		})

		return
	}
	LikedPost := &models.Like{}

	err = db.Find(&LikedPost, "post_id = ?", foundpost.ID).Error

	if err != nil && LikedPost.ID != 0 || LikedPost.UserID == user.ID {
		db.Exec("DELETE FROM likes where id=? ", LikedPost.ID)
		respondWithJSON(w, map[string]string{
			"msg": "Unliked post",
		})
		return
	}
	like := &models.Like{}

	like.UserID = user.ID
	like.PostID = foundpost.ID

	err = db.Create(&like).Error
	if err != nil {
		fmt.Println(err.Error())
		respondWithJSON(w, map[string]string{
			"error": "Couldnt like the post",
			"msg":   err.Error(),
		})

		return
	}
	// err = db.Preload("Likes").Preload("User").Find(&foundpost, "id = ?", request["id"]).Error
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	respondWithJSON(w, map[string]string{
	// 		"Error": "Post not found",
	// 		"Msg":   err.Error(),
	// 	})

	// 	return
	// }
	respondWithJSON(w, map[string]string{

		"msg": "Post Likes Successfully",
	})
}
