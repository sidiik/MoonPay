package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID        primitive.ObjectID   `bson:"_id"`
	UserID    int                  `bson:"user_id"`
	Balance   primitive.Decimal128 `bson:"balance"`
	Status    string               `bson:"status"`
	CreatedAt time.Time            `bson:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at"`
}

type WalletRepository interface {
	GetWalletByID(ctx context.Context, walletID string) (*Wallet, error)
	GetWalletByUserID(ctx context.Context, UserID int) (*Wallet, error)
	CreateWallet(ctx context.Context, userID int) (string, error)
}
