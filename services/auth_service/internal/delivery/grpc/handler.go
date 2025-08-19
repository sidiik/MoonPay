package handler

import (
	"context"

	"github.com/sidiik/moonpay/auth_service/internal/services"
	authpb "github.com/sidiik/moonpay/auth_service/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthServer struct {
	userService *services.UserService
	otpService  *services.OtpService
	authpb.UnimplementedAuthServiceServer
}

func NewAuthServerHandler(userService *services.UserService, otpService *services.OtpService) *AuthServer {
	return &AuthServer{
		userService: userService,
		otpService:  otpService,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	resp, err := s.userService.SignIn(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AuthServer) Signup(ctx context.Context, req *authpb.SignupRequest) (*authpb.SignupResponse, error) {
	user, err := s.userService.SignUp(ctx, req)
	if err != nil {
		return nil, err
	}

	return &authpb.SignupResponse{
		Email: user.Email,
	}, nil
}

func (s *AuthServer) RequestPasswordReset(ctx context.Context, req *authpb.RequestPasswordResetRequest) (*emptypb.Empty, error) {
	err := s.otpService.RequestPasswordReset(ctx, req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *AuthServer) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*emptypb.Empty, error) {
	err := s.otpService.ResetPassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
