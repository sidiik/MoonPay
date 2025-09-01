package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sidiik/moonpay/api_gateway/internal/constants"
	"github.com/sidiik/moonpay/api_gateway/internal/domain"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/authpb"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/wallet/walletpb"
	"github.com/sidiik/moonpay/api_gateway/internal/middleware"
	"github.com/sidiik/moonpay/api_gateway/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WalletHandler struct {
	walletServiceClient walletpb.WalletServiceClient
	authClient          authpb.AuthServiceClient
	validator           *validator.Validate
	log                 domain.Logger
}

func NewWalletHandler(log domain.Logger, walletServiceClient walletpb.WalletServiceClient, authClient authpb.AuthServiceClient, validator *validator.Validate, r *gin.RouterGroup) {
	handler := &WalletHandler{
		walletServiceClient: walletServiceClient,
		validator:           validator,
		log:                 log,
	}

	walletGroup := r.Group("/wallet")
	{
		walletGroup.POST("/request", middleware.Authenticate(), handler.RequestWallet)
		walletGroup.POST("/get-my-wallet", middleware.Authenticate(), handler.GetMyWallet)
	}

}

func (h *WalletHandler) RequestWallet(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
		c.Abort()
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
		c.Abort()
		return
	}

	md := metadata.New(map[string]string{
		"user-id": userIDStr,
	})

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	resp, err := h.walletServiceClient.RequestWallet(ctx, &emptypb.Empty{})

	if err != nil {
		s, _ := status.FromError(err)
		pkg.SendResponse(c, http.StatusBadRequest, s.Code().String(), s.Message(), nil, s.Err())
		c.Abort()
		return
	}

	pkg.SendResponse(c, http.StatusCreated, codes.OK.String(), "", resp, nil)

}

func (h *WalletHandler) GetMyWallet(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
		c.Abort()
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
		c.Abort()
		return
	}

	md := metadata.New(map[string]string{
		"user-id": userIDStr,
	})

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	resp, err := h.walletServiceClient.GetMyWallet(ctx, &emptypb.Empty{})

	if err != nil {
		s, _ := status.FromError(err)
		pkg.SendResponse(c, http.StatusBadRequest, s.Code().String(), s.Message(), nil, s.Err())
		c.Abort()
		return
	}

	pkg.SendResponse(c, http.StatusCreated, codes.OK.String(), "", resp, nil)

}
