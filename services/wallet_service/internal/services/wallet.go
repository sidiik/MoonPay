package services

import (
	"context"

	"github.com/sidiik/moonpay/wallet_service/internal/domain"
	"github.com/sidiik/moonpay/wallet_service/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletUsecase struct {
	WalletRepo domain.WalletRepository
	log        domain.Logger
}

func NewWalletUsecase(walletRepo domain.WalletRepository, log domain.Logger) *WalletUsecase {
	return &WalletUsecase{
		WalletRepo: walletRepo,
		log:        log,
	}
}

func (u *WalletUsecase) RequestWallet(ctx context.Context) (*string, error) {
	u.log.Info("checking if the user id is attached into the context")
	userID, err := pkg.GetUserIDFromCtx(ctx, u.log)
	if err != nil {
		return nil, err
	}

	if wallet, _ := u.GetMyWallet(ctx); wallet != nil {
		return nil, status.Error(codes.AlreadyExists, "wallet already exists")
	}

	u.log.Info("creating new wallet for user", "user-id", userID)

	walletID, err := u.WalletRepo.CreateWallet(ctx, *userID)
	if err != nil {
		u.log.Info("Unable to create wallet", "error", err)
		return nil, status.Error(codes.Internal, "unable to create wallet")
	}

	return &walletID, nil
}

func (u *WalletUsecase) GetMyWallet(ctx context.Context) (*domain.Wallet, error) {

	u.log.Info("checking if the user id is attached into the context")

	userID, err := pkg.GetUserIDFromCtx(ctx, u.log)
	if err != nil {
		return nil, err
	}

	wallet, err := u.WalletRepo.GetWalletByUserID(ctx, *userID)
	if err != nil {
		u.log.Error("unable to get user wallet", "error", err)
		return nil, status.Error(codes.NotFound, "user wallet not found")
	}

	return wallet, nil
}
