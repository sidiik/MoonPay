package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/authpb"
)

var validate *validator.Validate

func InitV1Routes(r *gin.Engine, client authpb.AuthServiceClient) {
	validate = validator.New(validator.WithRequiredStructEnabled())

	v1 := r.Group("/v1")
	{
		NewAuthHandler(client, validate, v1)
	}
}
