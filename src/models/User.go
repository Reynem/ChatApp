package models

import (
	"time"
)

type User struct {
	ID          uint       `json:"id"`
	UserName    string     `json:"user_name"`
	Email       string     `json:"email"`
	Password    string     `json:"password"`
	PhoneNumber string     `json:"phone_number"`
	BirthDate   *time.Time `json:"birth_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
