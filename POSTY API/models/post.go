package models

import (
	"time"
)

type Post struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	UserID    uint   `json:"userid" gorm:"column:user_id"`
	User      User
	Likes     []Like `json:"likes" gorm:"ForeignKey:PostID;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
