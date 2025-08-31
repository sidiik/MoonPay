package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        string          `bson:"_id"`
	UserID    int             `bson:"user_id"`
	Balance   decimal.Decimal `bson:"balance"`
	Status    string          `bson:"status"`
	CreatedAt time.Time       `bson:"created_at"`
	UpdatedAt time.Time       `bson:"updated_at"`
}

type WalletRepository interface {
	GetWalletByID(ctx context.Context, walletID string) (*Wallet, error)
	GetWalletByUserID(ctx context.Context, UserID int) (*Wallet, error)
	CreateWallet(ctx context.Context, userID int) (any, error)
}
