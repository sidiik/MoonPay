package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"userName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	RegisterUser(user *User) error
}

type UserService interface {
	RegisterUser(user *User) error
}
