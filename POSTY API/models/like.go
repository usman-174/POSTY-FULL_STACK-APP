package models

import (
	"time"
)

type Like struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	UserID    uint `gorm:"column:user_id" json:"userid"`
	User      User
	PostID    uint `gorm:"column:post_id" json:"postid"`
	Post      Post
	CreatedAt time.Time
	UpdatedAt time.Time
}
