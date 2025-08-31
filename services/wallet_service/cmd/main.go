package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "github.com/sidiik/moonpay/wallet_service/internal/delivery/grpc"
	"github.com/sidiik/moonpay/wallet_service/internal/infra/config"
	"github.com/sidiik/moonpay/wallet_service/internal/infra/logger"
	"github.com/sidiik/moonpay/wallet_service/internal/repository"
	"github.com/sidiik/moonpay/wallet_service/internal/services"
	walletpb "github.com/sidiik/moonpay/wallet_service/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger.Init()
	log := logger.NewZapLogger()
	config.InitConfig()
	appConfig := config.AppConfig

	log.Info("Initializing the wallet service")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.Port))

	if err != nil {
		log.Error("failed to listen wallet service", zap.Error(err))
		return
	}

	log.Info("connecting to mongodb")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(appConfig.MongoDBURI))
	if err != nil {
		log.Error("failed to connect to mongodb")
		return
	}

	db := client.Database("wallet_db")

	log.Info("initializing wallet repo")
	walletRepo := repository.NewWalletRepository(db)
	walletService := services.NewWalletUsecase(walletRepo, log)

	grpcServer := grpc.NewServer()

	walletpb.RegisterWalletServiceServer(grpcServer, handler.NewWalletServer(walletService, log))
	log.Info("wallet gRPC server running", "port", appConfig.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("failed to serve wallet service", zap.Error(err))
			return
		}
	}()

	<-stop

	done := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Info("wallet service stopped gracefully!")
	case <-time.After(10 * time.Second):
		log.Warn("Timeout reached, forcing server to stop")
		grpcServer.Stop()
	}

	log.Info("wallet service shutdown complete")
}
