package repository

import (
	"context"

	"github.com/sidiik/moonpay/auth_service/internal/domain"
	"github.com/sidiik/moonpay/auth_service/internal/infra/db"
)

type UserRepository struct {
	dbConnections *db.DBConnections
}

func NewUserRepository(dbConnections *db.DBConnections) *UserRepository {
	return &UserRepository{
		dbConnections: dbConnections,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	return r.dbConnections.Writer.WithContext(ctx).Create(&user).Error
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	if err := r.dbConnections.Reader.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	if err := r.dbConnections.Reader.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ResetPassword(ctx context.Context, userID uint, newPassword string) error {
	return r.dbConnections.Writer.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Update("password", newPassword).Error
}
