package wallet

import (
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/wallet/walletpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewWalletServiceClient(conn *grpc.ClientConn) walletpb.WalletServiceClient {
	return walletpb.NewWalletServiceClient(conn)
}

func ConnectWalletService(address string) (walletpb.WalletServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return NewWalletServiceClient(conn), nil
}
