package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string
	FullName     string
	Email        string `gorm:"unqiueIdx"`
	Is2faEnabled bool
	IsLocked     bool
	LockUntil    time.Time
	Password     string
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	ResetPassword(ctx context.Context, userID uint, newPassword string) error
}
