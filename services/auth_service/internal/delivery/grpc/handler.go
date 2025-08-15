package handler

import (
	"context"

	"github.com/sidiik/moonpay/auth_service/internal/services"
	authpb "github.com/sidiik/moonpay/auth_service/proto"
)

type AuthServer struct {
	service *services.AuthService
	authpb.UnimplementedAuthServiceServer
}

func NewAuthServerHandler(svc *services.AuthService) *AuthServer {
	return &AuthServer{
		service: svc,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	resp, err := s.service.SignIn(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AuthServer) Signup(ctx context.Context, req *authpb.SignupRequest) (*authpb.SignupResponse, error) {
	user, err := s.service.SignUp(ctx, req)
	if err != nil {
		return nil, err
	}

	return &authpb.SignupResponse{
		Email: user.Email,
	}, nil
}
