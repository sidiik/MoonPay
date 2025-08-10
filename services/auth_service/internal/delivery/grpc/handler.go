package grpc_handler

import (
	"context"

	auth "github.com/sidiik/moonpay/auth_service/internal/genproto/proto"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
}

func NewAuthServer() *AuthServer {
	return &AuthServer{}
}

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	return &auth.LoginResponse{
		AccessToken:  req.Email,
		RefreshToken: req.Password,
	}, nil
}
