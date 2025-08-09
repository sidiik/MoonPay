package grpc

import (
	"context"

	pb "github.com/sidiik/moonpay/auth_service/internal/delivery/grpc/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthServer() *AuthServer {
	return &AuthServer{}
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}, nil
}
