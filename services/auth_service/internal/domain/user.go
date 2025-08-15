package domain

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	FullName string
	Email    string `gorm:"unqiueIdx"`
	Password string
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
}
