package main

import (
	"fmt"
	"log/slog"
	"net"

	grpc_handler "github.com/sidiik/moonpay/auth_service/internal/delivery/grpc"
	auth "github.com/sidiik/moonpay/auth_service/internal/genproto/proto"
	"github.com/sidiik/moonpay/auth_service/internal/infra/config"
	"google.golang.org/grpc"
)

func main() {
	slog.Info("Initializing auth service env variables")
	config.InitConfig()

	slog.Info("Starting auth service")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.AppConfig.Port))
	if err != nil {
		slog.Error("failed to start auth service", "error", err)
		return
	}

	grpcServer := grpc.NewServer()

	// Registering the auth service
	auth.RegisterAuthServiceServer(grpcServer, grpc_handler.NewAuthServer())
	slog.Info(fmt.Sprintf("🚀 gRPC AuthService running on :%s", config.AppConfig.Port))

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve auth server", "error", err)
		return
	}
}
