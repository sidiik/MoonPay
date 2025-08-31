package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "github.com/sidiik/moonpay/auth_service/internal/delivery/grpc"
	"github.com/sidiik/moonpay/auth_service/internal/infra/config"
	"github.com/sidiik/moonpay/auth_service/internal/infra/db"
	"github.com/sidiik/moonpay/auth_service/internal/infra/logger"
	"github.com/sidiik/moonpay/auth_service/internal/infra/rabbitmq"
	"github.com/sidiik/moonpay/auth_service/internal/infra/redis"
	"github.com/sidiik/moonpay/auth_service/internal/repository"
	"github.com/sidiik/moonpay/auth_service/internal/services"
	authpb "github.com/sidiik/moonpay/auth_service/proto"
	"google.golang.org/grpc"
)

func main() {
	// Initializing global logger
	logger.Init()
	log := logger.NewZapLogger()

	log.Info("Initializing auth service env variables")
	config.InitConfig()
	appConfig := config.AppConfig

	log.Info("Initializing redis db")
	redisClient := redis.InitClient(appConfig)

	log.Info("Initializing DB Connections")
	conns, err := db.InitDb()
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		return
	}

	// Initialize rabbitmq
	r, err := rabbitmq.NewRabbitMQ(appConfig.RabbitMQUrl)
	if err != nil {
		log.Error("failed to connect to rabbitmq", "error", err)
		return
	}

	defer r.Close()

	// Initialize auth service and repo
	userRepo := repository.NewUserRepository(conns)
	userService := services.NewUserService(userRepo, r, log, redisClient)

	// Initialize otp service and repo
	otpRepo := repository.NewOtpRepository(conns)
	otpService := services.NewOtpService(otpRepo, userRepo, r, log)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.Port))
	if err != nil {
		log.Error("failed to start auth service", "error", err)
		return
	}

	grpcServer := grpc.NewServer()

	// Registering the auth service
	authpb.RegisterAuthServiceServer(grpcServer, handler.NewAuthServerHandler(userService, otpService))
	log.Info(fmt.Sprintf("🚀 gRPC AuthService running on :%s", appConfig.Port))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("failed to serve auth server", "error", err)
			return
		}
	}()

	<-stop
	log.Info("Shutting down auth gRPC server gracefully...")

	done := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Info("auth gRPC server stopped gracefully")
	case <-time.After(10 * time.Second):
		log.Warn("Timeout reached, forcing server stop")
		grpcServer.Stop()
	}

	log.Info("Auth service shutdown complete")
}
