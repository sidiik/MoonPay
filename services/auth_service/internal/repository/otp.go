package repository

import (
	"context"
	"time"

	"github.com/sidiik/moonpay/auth_service/internal/domain"
	"github.com/sidiik/moonpay/auth_service/internal/infra/db"
)

type OtpRepository struct {
	dbConnections *db.DBConnections
}

func NewOtpRepository(dbConnections *db.DBConnections) *OtpRepository {
	return &OtpRepository{
		dbConnections: dbConnections,
	}
}

func (r *OtpRepository) CheckActiveOtp(ctx context.Context, userID uint) (*domain.Otp, error) {
	var otp domain.Otp

	if err := r.dbConnections.
		Reader.
		WithContext(ctx).
		Where("user_id = ? AND expires_at > ? AND is_used = ?", userID, time.Now().UTC(), false).
		Order("created_at DESC").
		First(&otp).
		Error; err != nil {
		return nil, err
	}

	return &otp, nil
}

func (r *OtpRepository) CreateOtp(ctx context.Context, otp *domain.Otp) error {
	if err := r.dbConnections.Writer.Create(&otp).Error; err != nil {
		return err
	}

	return nil
}

func (r *OtpRepository) MarkOtpAsUsed(ctx context.Context, otpID uint) error {
	if err := r.dbConnections.Writer.WithContext(ctx).Where("id = ?", otpID).Model(&domain.Otp{}).Update("is_used", true).Error; err != nil {
		return err
	}

	return nil
}
