package repository

import (
	"context"
	"time"

	"github.com/sidiik/moonpay/wallet_service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *WalletRepository) CreateWallet(ctx context.Context, userID int) (string, error) {
	balance, _ := primitive.ParseDecimal128("0")
	wallet := &domain.Wallet{
		ID:        primitive.NewObjectID(),
		Balance:   balance,
		UserID:    userID,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, wallet)
	return wallet.ID.Hex(), err
}
