package handler

import (
	"context"

	walletpb "github.com/sidiik/moonpay/wallet_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletServer struct {
	walletpb.UnimplementedWalletServiceServer
}

func NewWalletServer() *WalletServer {
	return &WalletServer{}
}

func (w *WalletServer) RequestWallet(ctx context.Context, req *walletpb.RequestWalletRequest) (*walletpb.RequestWalletResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Request wallet is not implemented yet!")
}
