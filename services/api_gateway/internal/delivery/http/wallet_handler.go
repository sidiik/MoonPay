package http

import (
	"net/http"
	"strconv"

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
		authClient:          authClient,
		log:                 log,
	}

	walletGroup := r.Group("/wallet")
	{
		walletGroup.POST("/request", middleware.Authenticate(handler.authClient), handler.RequestWallet)
	}

}

func (h *WalletHandler) RequestWallet(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
		c.Abort()
		return
	}

	userd := user.(*authpb.GetUserByEmailResponse)

	userIDStr := strconv.Itoa(int(userd.Id))

	md := metadata.New(map[string]string{
		"user-id": userIDStr,
	})

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	resp, err := h.walletServiceClient.RequestWallet(ctx, &walletpb.RequestWalletRequest{
		UserId: strconv.Itoa(int(userd.Id)),
	})

	if err != nil {
		s, _ := status.FromError(err)
		pkg.SendResponse(c, http.StatusBadRequest, s.Code().String(), s.Message(), nil, s.Err())
		c.Abort()
		return
	}

	pkg.SendResponse(c, http.StatusCreated, codes.OK.String(), "", resp, nil)

}
