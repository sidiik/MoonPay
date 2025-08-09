package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sidiik/moonpay/api_gateway/internal/constants"
	"github.com/sidiik/moonpay/api_gateway/internal/dto"
	pb "github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/proto"
	"github.com/sidiik/moonpay/api_gateway/pkg"
)

type AuthHandler struct {
	authClient pb.AuthServiceClient
	validator  *validator.Validate
}

func NewAuthHandler(client pb.AuthServiceClient, validator *validator.Validate, r *gin.RouterGroup) {

	handler := &AuthHandler{
		authClient: client,
		validator:  validator,
	}

	auth := r.Group("/auth")
	{
		auth.POST("/login", handler.Login)
	}

}

func (h *AuthHandler) Login(c *gin.Context) {
	var body dto.LoginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Warn("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, constants.ErrInvalidRequest, nil)
		return
	}

	ctx := c.Request.Context()

	resp, err := h.authClient.Login(ctx, &pb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		slog.Warn("error from auth grpc service", "error", err)
		pkg.SendResponse(c, http.StatusUnauthorized, constants.ErrUnauthorized, nil)
		return
	}

	pkg.SendResponse(c, http.StatusOK, "", resp)

}
