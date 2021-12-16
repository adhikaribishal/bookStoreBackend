package models

import "time"

type User struct {
	ID        string    `json:"id" sql:"id"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	Username  string    `json:"username" sql:"username"`
	FirstName string    `json:"first_name" sql:"first_name"`
	LastName  string    `json:"last_name" sql:"last_name"`
	CreatedAt time.Time `json:"createdat" sql:"created_at"`
	UpdatedAt time.Time `json:"updatedat" sql:"updated_at"`
}
