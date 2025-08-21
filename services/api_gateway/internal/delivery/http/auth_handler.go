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
	"github.com/sidiik/moonpay/api_gateway/internal/middleware"
	"github.com/sidiik/moonpay/api_gateway/pkg"
	"google.golang.org/grpc/codes"
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
		auth.POST("/request-password-reset", handler.RequestPasswordReset)
		auth.POST("/reset-password", handler.ResetPassword)
		auth.GET("/whoami", middleware.Authenticate(handler.authClient), handler.WhoAmI)
	}

}

func (h *AuthHandler) Signup(c *gin.Context) {
	var body dto.SignupRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Error("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, constants.ErrInvalidRequest, "", nil, err)
		return
	}

	if err := h.validator.Struct(&body); err != nil {
		slog.Error("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, constants.ErrInvalidRequest, "", nil, err)
		return
	}

	ctx := c.Request.Context()

	resp, err := h.authClient.Signup(ctx, &authpb.SignupRequest{
		Email:    body.Email,
		Password: body.Password,
		FullName: body.FullName,
	})

	s, _ := status.FromError(err)

	if err != nil {
		slog.Error("error from auth grpc service", "error", err)
		pkg.SendResponse(c, http.StatusUnauthorized, s.Code().String(), "", nil, errors.New(s.Message()))
		return
	}

	pkg.SendResponse(c, http.StatusOK, "", "", resp, nil)

}

func (h *AuthHandler) Signin(c *gin.Context) {
	var body dto.SigninRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Error("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, "", "", nil, err)
		return
	}

	if err := h.validator.Struct(&body); err != nil {
		slog.Error("failed to parse body", "error", err)
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
		slog.Error("error from auth grpc service", "error", err)
		pkg.SendResponse(c, http.StatusUnauthorized, s.Code().String(), "", nil, errors.New(s.Message()))
		return
	}

	pkg.SendResponse(c, http.StatusOK, "", "", resp, nil)

}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var body dto.RequestResetPasswordRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Error("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, "", "", nil, err)
		return
	}

	if err := h.validator.Struct(&body); err != nil {
		slog.Error("failed to validate body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, "", "", nil, err)
		return
	}

	ctx := c.Request.Context()

	resp, err := h.authClient.RequestPasswordReset(ctx, &authpb.RequestPasswordResetRequest{
		Email: body.Email,
	})

	s, _ := status.FromError(err)

	if err != nil {
		slog.Error("failed to request password reset", "error", err)
		pkg.SendResponse(c, http.StatusUnauthorized, s.Code().String(), "", nil, errors.New(s.Message()))
		return
	}

	pkg.SendResponse(c, http.StatusOK, "", "", resp, nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var body dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Error("failed to parse body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, "", "", nil, err)
		return
	}

	if err := h.validator.Struct(&body); err != nil {
		slog.Error("failed to validate body", "error", err)
		pkg.SendResponse(c, http.StatusBadRequest, "", "", nil, err)
		return
	}

	ctx := c.Request.Context()

	resp, err := h.authClient.ResetPassword(ctx, &authpb.ResetPasswordRequest{
		Email:       body.Email,
		Code:        body.Code,
		NewPassword: body.NewPassword,
	})

	s, _ := status.FromError(err)

	if err != nil {
		slog.Error("failed to request password reset", "error", err)
		pkg.SendResponse(c, http.StatusUnauthorized, s.Code().String(), "", nil, errors.New(s.Message()))
		return
	}

	pkg.SendResponse(c, http.StatusOK, "", "", resp, nil)
}

func (h *AuthHandler) WhoAmI(c *gin.Context) {
	user := c.MustGet("user").(*authpb.GetUserByEmailResponse)
	pkg.SendResponse(c, http.StatusOK, codes.OK.String(), "", user, nil)
}
