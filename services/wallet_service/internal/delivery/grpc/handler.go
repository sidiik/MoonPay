package handler

import (
	"context"

	"github.com/sidiik/moonpay/wallet_service/internal/domain"
	"github.com/sidiik/moonpay/wallet_service/internal/services"
	walletpb "github.com/sidiik/moonpay/wallet_service/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (w *WalletServer) RequestWallet(ctx context.Context, req *emptypb.Empty) (*walletpb.RequestWalletResponse, error) {
	w.log.Info("Creating new wallet for user")
	walletID, err := w.walletService.RequestWallet(ctx)
	if err != nil {
		return nil, err
	}

	return &walletpb.RequestWalletResponse{
		WalletId: *walletID,
	}, nil
}

func (w *WalletServer) GetMyWallet(ctx context.Context, req *emptypb.Empty) (*walletpb.GetMyWalletResponse, error) {
	w.log.Info("Getting my wallet for user")
	wallet, err := w.walletService.GetMyWallet(ctx)
	if err != nil {
		return nil, err
	}

	return &walletpb.GetMyWalletResponse{
		WalletId:  wallet.ID.Hex(),
		Balance:   wallet.Balance.String(),
		UserId:    int64(wallet.UserID),
		CreatedAt: timestamppb.New(wallet.CreatedAt),
		UpdatedAt: timestamppb.New(wallet.UpdatedAt),
	}, nil
}
