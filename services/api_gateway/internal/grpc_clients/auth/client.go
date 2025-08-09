package auth

import (
	"log/slog"

	pb "github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(conn *grpc.ClientConn) pb.AuthServiceClient {
	return pb.NewAuthServiceClient(conn)
}

func ConnectAuthService(address string) pb.AuthServiceClient {
	slog.Info("Connecting auth service grpc client")
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to connect to auth_service", "error", err)
	}

	return NewAuthClient(conn)
}
