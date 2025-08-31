package handler

import (
	"context"

	"github.com/sidiik/moonpay/wallet_service/internal/domain"
	"github.com/sidiik/moonpay/wallet_service/internal/services"
	walletpb "github.com/sidiik/moonpay/wallet_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletServer struct {
	walletService *services.WalletUsecase
	log           domain.Logger
	walletpb.UnimplementedWalletServiceServer
}

func NewWalletServer(walletService *services.WalletUsecase, log domain.Logger) *WalletServer {
	return &WalletServer{
		walletService: walletService,
		log:           log,
	}
}

func (w *WalletServer) RequestWallet(ctx context.Context, req *walletpb.RequestWalletRequest) (*walletpb.RequestWalletResponse, error) {
	w.log.Info("Creating new wallet for user")
	walletID, err := w.walletService.RequestWallet(ctx)
	if err != nil {
		return nil, err
	}

	walletIdStr, ok := walletID.(string)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to convert wallet id to string")
	}

	return &walletpb.RequestWalletResponse{
		WalletId: walletIdStr,
	}, nil
}
