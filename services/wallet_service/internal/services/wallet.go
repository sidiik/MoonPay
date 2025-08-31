package services

import "github.com/sidiik/moonpay/wallet_service/internal/domain"

type WalletUsecase struct {
	WalletRepo domain.WalletRepository
}

func NewWalletRepository(walletRepo domain.WalletRepository) *WalletUsecase {

}
