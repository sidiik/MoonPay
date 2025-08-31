package repository

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/sidiik/moonpay/wallet_service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WalletRepository struct {
	collection *mongo.Collection
}

func NewWalletRepository(db *mongo.Database) domain.WalletRepository {
	return &WalletRepository{
		collection: db.Collection("wallets"),
	}
}

func (r *WalletRepository) GetWalletByID(ctx context.Context, walletID string) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.collection.FindOne(ctx, bson.M{"_id": walletID}).Decode(&wallet)
	return &wallet, err
}

func (r *WalletRepository) GetWalletByUserID(ctx context.Context, UserID int) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.collection.FindOne(ctx, bson.M{"user_id": UserID}).Decode(&wallet)
	return &wallet, err
}

func (r *WalletRepository) CreateWallet(ctx context.Context, userID int) (any, error) {
	wallet := &domain.Wallet{
		Balance: decimal.NewFromInt(0),
		UserID:  userID,
	}

	result, err := r.collection.InsertOne(ctx, wallet)
	return result.InsertedID, err
}
