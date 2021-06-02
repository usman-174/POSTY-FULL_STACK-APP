package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Username  string `json:"username"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Posts     []Post `json:"posts" gorm:"ForeignKey:UserID;constraint:OnDelete:CASCADE;"`
}
