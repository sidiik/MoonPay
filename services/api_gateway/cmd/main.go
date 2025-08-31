package main

import (
	"github.com/gin-gonic/gin"
	h "github.com/sidiik/moonpay/api_gateway/internal/delivery/http"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/wallet"
	"github.com/sidiik/moonpay/api_gateway/internal/infra/config"
	"github.com/sidiik/moonpay/api_gateway/internal/infra/logger"
)

func main() {
	logger.Init()
	log := logger.NewZapLogger()

	log.Info("Initializing the config")
	config.InitConfig()

	log.Info("Starting API Gateway service")
	r := gin.Default()

	// Init auth service
	log.Info("Initializing auth service gRPC")
	authServiceClient, err := auth.ConnectAuthService("auth-service:80")
	if err != nil {
		log.Error("failed to start auth service client", "error", err)
		return
	}
	walletServiceClient, err := wallet.ConnectWalletService("wallet-service:80")
	if err != nil {
		log.Error("failed to start auth service client", "error", err)
		return
	}

	h.InitV1Routes(r, log, authServiceClient, walletServiceClient)

	// Listen
	r.Run(":8080")

}
