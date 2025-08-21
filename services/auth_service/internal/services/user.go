package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sidiik/moonpay/auth_service/internal/constants"
	"github.com/sidiik/moonpay/auth_service/internal/domain"
	"github.com/sidiik/moonpay/auth_service/internal/infra/rabbitmq"
	"github.com/sidiik/moonpay/auth_service/internal/infra/redis"
	"github.com/sidiik/moonpay/auth_service/internal/utils"
	authpb "github.com/sidiik/moonpay/auth_service/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	userRepo    domain.UserRepository
	rabbit      *rabbitmq.RabbitMQ
	log         domain.Logger
	redisClient *redis.RedisClient
}

func NewUserService(userRepo domain.UserRepository, rabbit *rabbitmq.RabbitMQ, log domain.Logger, redisClient *redis.RedisClient) *UserService {
	return &UserService{
		userRepo:    userRepo,
		rabbit:      rabbit,
		log:         log,
		redisClient: redisClient,
	}
}

func (s *UserService) SignUp(ctx context.Context, data *authpb.SignupRequest) (*domain.User, error) {
	data.Email = strings.Trim(strings.ToLower(data.Email), " ")
	if existingUser, _ := s.userRepo.GetUserByEmail(ctx, data.Email); existingUser != nil {
		s.log.Error("failed to create user", "email", data.Email, "error", constants.ErrEmailAlreadyExists)
		return nil, status.Error(codes.AlreadyExists, constants.ErrEmailAlreadyExists)
	}

	s.log.Info("creating new user", "email", data.Email)

	s.log.Info("Hashing the password")
	hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if err != nil {
		s.log.Error("failed to create user", "email", data.Email, "error", err)
		return nil, status.Error(codes.Internal, constants.ErrInternalServer)
	}

	user := domain.User{
		Username: strings.Split(data.Email, "@")[0],
		Email:    data.Email,
		FullName: data.FullName,
		Password: string(hashed),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		s.log.Error("failed to create user", "email", data.Email, "error", err)
	}

	event := map[string]any{
		"event":    "user.registered",
		"email":    user.Email,
		"fullName": user.FullName,
	}

	body, err := json.Marshal(event)
	if err != nil {
		s.log.Error("failed to marshal rabbitmq event")
	}

	s.log.Info("Publish user.registered event")
	if err := s.rabbit.Publish("user.registered", body); err != nil {
		s.log.Error("failed to publish event", "error", err)
	}

	return &user, nil

}

func (s *UserService) SignIn(ctx context.Context, data *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	data.Email = strings.Trim(strings.ToLower(data.Email), " ")
	existingUser, err := s.userRepo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		s.log.Error("failed to get user by email", "email", data.Email, "error", constants.ErrInvalidCredentials)
		return nil, status.Error(codes.Unauthenticated, constants.ErrInvalidCredentials)
	}

	s.log.Info("Check if the password is correct")
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(data.Password)); err != nil {
		s.log.Error("incorrect password", "email", data.Email, "error", err)
		return nil, status.Error(codes.Unauthenticated, constants.ErrInvalidCredentials)
	}

	accessToken, err := utils.GenerateAccessToken(existingUser.Email)
	if err != nil {
		s.log.Error("failed to generate access token", "error", err)
		return nil, status.Error(codes.Internal, constants.ErrInternalServer)
	}

	refreshToken, err := utils.GenerateRefreshToken(existingUser.Email)
	if err != nil {
		s.log.Error("failed to generate refresh token", "error", err)
		return nil, status.Error(codes.Internal, constants.ErrInternalServer)
	}

	return &authpb.LoginResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil

}

func (s *UserService) GetUserByEmail(ctx context.Context, req *authpb.GetUserByEmailRequest) (*authpb.GetUserByEmailResponse, error) {
	s.log.Info("checking if the user is stored in redis")
	userStr, err := s.redisClient.Get(ctx, fmt.Sprintf("user::%s", req.Email))

	if err != nil {
		s.log.Warn("failed to get user from redis", "error", err)
	}

	if userStr != "" {
		var existingUser domain.User
		err := json.Unmarshal([]byte(userStr), &existingUser)
		if err == nil {
			return &authpb.GetUserByEmailResponse{
				Id:        uint64(existingUser.ID),
				Email:     existingUser.Email,
				FullName:  existingUser.FullName,
				CreatedAt: timestamppb.New(existingUser.CreatedAt),
				UpdatedAt: timestamppb.New(existingUser.UpdatedAt),
			}, nil
		}
	}

	s.log.Info("checking if the user exists from the database")
	req.Email = strings.Trim(strings.ToLower(req.Email), " ")
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.log.Error("failed to get user by email", "email", req.Email, "error", constants.ErrInvalidCredentials)
		return nil, status.Error(codes.Unauthenticated, constants.ErrInvalidCredentials)
	}

	s.log.Info("Setting password as empty before storing into redis")
	existingUser.Password = ""

	body, err := json.Marshal(existingUser)
	if err != nil {
		s.log.Error("failed to marshal the user before storing into redis")
	}

	s.log.Info("storing user into redis")
	exp := time.Hour * time.Duration(1)
	if err := s.redisClient.Set(ctx, exp, fmt.Sprintf("user::%s", existingUser.Email), string(body)); err != nil {
		s.log.Error("failed to store the user in redis", "error", err)
	}

	return &authpb.GetUserByEmailResponse{
		Id:        uint64(existingUser.ID),
		Email:     existingUser.Email,
		FullName:  existingUser.FullName,
		CreatedAt: timestamppb.New(existingUser.CreatedAt),
		UpdatedAt: timestamppb.New(existingUser.UpdatedAt),
	}, nil

}
