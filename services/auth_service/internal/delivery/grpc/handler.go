package handler

import (
	"context"

	authpb "github.com/sidiik/moonpay/auth_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthServer() *AuthServer {
	return &AuthServer{}
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signup not implemented")
}

func (s *AuthServer) Signup(context.Context, *authpb.SignupRequest) (*authpb.SignupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signup not implemented")
}
