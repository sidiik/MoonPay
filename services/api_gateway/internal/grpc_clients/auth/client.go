package auth

import (
	"log/slog"

	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(conn *grpc.ClientConn) authpb.AuthServiceClient {
	return authpb.NewAuthServiceClient(conn)
}

func ConnectAuthService(address string) authpb.AuthServiceClient {
	slog.Info("Connecting auth service grpc client")
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to connect to auth_service", "error", err)
	}

	return NewAuthClient(conn)
}
