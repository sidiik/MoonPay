package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	PasswordReset = iota + 1
	TwoFactorAuthentication
)

type Otp struct {
	gorm.Model
	Code         string
	ExpiresAt    time.Time
	NextResendAt time.Time
	OtpReason    int
	IsUsed       bool
	UsedAt       time.Time
	UserID       uint
	User         User `gorm:"foreignKey:UserID"`
}

type OtpRepository interface {
	CheckActiveOtp(ctx context.Context, userID uint) (*Otp, error)
	CreateOtp(ctx context.Context, otp *Otp) error
	MarkOtpAsUsed(ctx context.Context, otpID uint) error
}
