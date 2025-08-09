package main

import (
	"log/slog"
	"net"

	auth_grpc "github.com/sidiik/moonpay/auth_service/internal/delivery/grpc"
	pb "github.com/sidiik/moonpay/auth_service/internal/delivery/grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	slog.Info("Starting auth service")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("failed to start auth service", "error", err)
		return
	}

	grpcServer := grpc.NewServer()

	// Registering the auth service
	pb.RegisterAuthServiceServer(grpcServer, auth_grpc.NewAuthServer())
	slog.Info("🚀 gRPC AuthService running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve auth server", "error", err)
		return
	}
}
