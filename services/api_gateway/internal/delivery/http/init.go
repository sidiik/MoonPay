package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sidiik/moonpay/api_gateway/internal/domain"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/authpb"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/wallet/walletpb"
)

var validate *validator.Validate

func InitV1Routes(r *gin.Engine, log domain.Logger, authServiceClient authpb.AuthServiceClient, walletServiceClient walletpb.WalletServiceClient) {
	validate = validator.New(validator.WithRequiredStructEnabled())

	v1 := r.Group("/v1")
	{
		NewAuthHandler(authServiceClient, validate, v1)
		NewWalletHandler(log, walletServiceClient, authServiceClient, validate, v1)
	}
}
