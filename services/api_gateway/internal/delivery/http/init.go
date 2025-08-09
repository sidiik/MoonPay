package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	pb "github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/proto"
)

var validate *validator.Validate

func InitV1Routes(r *gin.Engine, client pb.AuthServiceClient) {
	validate = validator.New(validator.WithRequiredStructEnabled())

	v1 := r.Group("/v1")
	{
		NewAuthHandler(client, validate, v1)
	}
}
