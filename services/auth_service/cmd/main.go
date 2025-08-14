package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "github.com/sidiik/moonpay/auth_service/internal/delivery/grpc"
	"github.com/sidiik/moonpay/auth_service/internal/infra/config"
	authpb "github.com/sidiik/moonpay/auth_service/proto"
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
	authpb.RegisterAuthServiceServer(grpcServer, handler.NewAuthServer())
	slog.Info(fmt.Sprintf("🚀 gRPC AuthService running on :%s", config.AppConfig.Port))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("failed to serve auth server", "error", err)
			return
		}
	}()

	<-stop
	slog.Info("Shutting down auth gRPC server gracefully...")

	done := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("auth gRPC server stopped gracefully")
	case <-time.After(10 * time.Second):
		slog.Warn("Timeout reached, forcing server stop")
		grpcServer.Stop()
	}

	slog.Info("Auth service shutdown complete")
}
