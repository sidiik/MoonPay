package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sidiik/moonpay/api_gateway/internal/constants"
	"github.com/sidiik/moonpay/api_gateway/internal/dto"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/authpb"
	"github.com/sidiik/moonpay/api_gateway/pkg"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authClient authpb.AuthServiceClient
	validator  *validator.Validate
}

func NewAuthHandler(client authpb.AuthServiceClient, validator *validator.Validate, r *gin.RouterGroup) {

	handler := &AuthHandler{
		authClient: client,
		validator:  validator,
	}

	auth := r.Group("/auth")
	{
		auth.POST("/signup", handler.Signup)
		auth.POST("/signin", handler.Signin)
	}

}

func (h *AuthHandler) Signup(c *gin.Context) {
	var body dto.SignupRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Warn("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, constants.ErrInvalidRequest, "", nil, err)
		return
	}

	if err := h.validator.Struct(&body); err != nil {
		slog.Warn("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, constants.ErrInvalidRequest, "", nil, err)
		return
	}

	ctx := c.Request.Context()

	resp, err := h.authClient.Signup(ctx, &authpb.SignupRequest{
		Email:    body.Email,
		Password: body.Password,
		FullName: body.FullName,
	})

	if err != nil {
		slog.Warn("error from auth grpc service", "error", err)
		pkg.SendResponse(c, http.StatusUnauthorized, constants.ErrUnauthorized, "", nil, errors.New(constants.ErrInternal))
		return
	}

	pkg.SendResponse(c, http.StatusOK, "", "", resp, nil)

}

func (h *AuthHandler) Signin(c *gin.Context) {
	var body dto.SigninRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Warn("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, "", "", nil, err)
		return
	}

	if err := h.validator.Struct(&body); err != nil {
		slog.Warn("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, "", "", nil, err)
		return
	}

	ctx := c.Request.Context()

	resp, err := h.authClient.Login(ctx, &authpb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	s, _ := status.FromError(err)

	if err != nil {
		slog.Warn("error from auth grpc service", "error", err)
		pkg.SendResponse(c, http.StatusUnauthorized, s.Code().String(), "", nil, errors.New(s.Message()))
		return
	}

	pkg.SendResponse(c, http.StatusOK, "", "", resp, nil)

}
