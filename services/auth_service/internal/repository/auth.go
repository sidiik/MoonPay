package repository

import (
	"context"

	"github.com/sidiik/moonpay/auth_service/internal/domain"
	"github.com/sidiik/moonpay/auth_service/internal/infra/db"
)

type AuthRepository struct {
	dbConnections *db.DBConnections
}

func NewAuthRepository(dbConnections *db.DBConnections) *AuthRepository {
	return &AuthRepository{
		dbConnections: dbConnections,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user domain.User) error {
	return r.dbConnections.Writer.Create(&user).Error
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	if err := r.dbConnections.Reader.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	if err := r.dbConnections.Reader.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
