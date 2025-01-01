package models

import "time"

type User struct {
	UserID    int       `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"size:50"`
	Email     string    `json:"email" gorm:"size:100;unique"`
	Password  string    `json:"password" gorm:"column:password_hash;size:255"`
	Role      string    `json:"role" gorm:"size:20"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
