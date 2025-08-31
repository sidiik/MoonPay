package auth

import (
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(conn *grpc.ClientConn) authpb.AuthServiceClient {
	return authpb.NewAuthServiceClient(conn)
}

func ConnectAuthService(address string) (authpb.AuthServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return NewAuthClient(conn), nil
}
