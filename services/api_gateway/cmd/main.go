package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	h "github.com/sidiik/moonpay/api_gateway/internal/delivery/http"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth"
)

func main() {
	slog.Info("Starting API Gateway")
	r := gin.Default()

	// Init auth service
	authServiceClient := auth.ConnectAuthService("localhost:50051")

	h.InitV1Routes(r, authServiceClient)

	// Listen
	r.Run(":8080")

}
