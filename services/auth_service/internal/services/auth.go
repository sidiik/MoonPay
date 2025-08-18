package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/sidiik/moonpay/auth_service/internal/constants"
	"github.com/sidiik/moonpay/auth_service/internal/domain"
	"github.com/sidiik/moonpay/auth_service/internal/infra/rabbitmq"
	"github.com/sidiik/moonpay/auth_service/internal/repository"
	"github.com/sidiik/moonpay/auth_service/internal/utils"
	authpb "github.com/sidiik/moonpay/auth_service/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	authRepo *repository.AuthRepository
	rabbit   *rabbitmq.RabbitMQ
}

func NewAuthService(authRepo *repository.AuthRepository, rabbit *rabbitmq.RabbitMQ) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		rabbit:   rabbit,
	}
}

func (s *AuthService) SignUp(ctx context.Context, data *authpb.SignupRequest) (*domain.User, error) {
	data.Email = strings.Trim(strings.ToLower(data.Email), " ")
	if existingUser, _ := s.authRepo.GetUserByEmail(ctx, data.Email); existingUser != nil {
		slog.Error("failed to create user", "email", data.Email, "error", constants.ErrEmailAlreadyExists)
		return nil, status.Error(codes.AlreadyExists, constants.ErrEmailAlreadyExists)
	}

	slog.Info("creating new user", "email", data.Email)

	slog.Info("Hashing the password")
	hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("failed to create user", "email", data.Email, "error", err)
		return nil, status.Error(codes.Internal, constants.ErrInternalServer)
	}

	user := domain.User{
		Username: strings.Split(data.Email, "@")[0],
		Email:    data.Email,
		FullName: data.FullName,
		Password: string(hashed),
	}

	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		slog.Error("failed to create user", "email", data.Email, "error", err)
	}

	event := map[string]any{
		"event":     "user.registered",
		"email":     user.Email,
		"full_name": user.FullName,
	}

	body, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal rabbitmq event")
	}

	slog.Info("Publish user.registered event")
	if err := s.rabbit.Publish("user.registered", body); err != nil {
		slog.Error("failed to publish event", "error", err)
	}

	return &user, nil

}

func (s *AuthService) SignIn(ctx context.Context, data *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	data.Email = strings.Trim(strings.ToLower(data.Email), " ")
	existingUser, err := s.authRepo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		slog.Error("failed to get user by email", "email", data.Email, "error", constants.ErrInvalidCredentials)
		return nil, status.Error(codes.Unauthenticated, constants.ErrInvalidCredentials)
	}

	slog.Info("Check if the password is correct")
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(data.Password)); err != nil {
		slog.Error("incorrect password", "email", data.Email, "error", err)
		return nil, status.Error(codes.Unauthenticated, constants.ErrInvalidCredentials)
	}

	accessToken, err := utils.GenerateAccessToken(existingUser.Email)
	if err != nil {
		slog.Error("failed to generate access token", "error", err)
		return nil, status.Error(codes.Internal, constants.ErrInternalServer)
	}

	refreshToken, err := utils.GenerateRefreshToken(existingUser.Email)
	if err != nil {
		slog.Error("failed to generate refresh token", "error", err)
		return nil, status.Error(codes.Internal, constants.ErrInternalServer)
	}

	return &authpb.LoginResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil

}
