package pkg

import (
	"context"
	"strconv"

	"github.com/sidiik/moonpay/wallet_service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GetUserIDFromCtx(ctx context.Context, log domain.Logger) (*int, error) {
	userIDs := metadata.ValueFromIncomingContext(ctx, "user-id")
	if len(userIDs) == 0 {
		log.Warn("user ids not found in the ctx")
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	userIDStr := userIDs[0]
	if userIDStr == "" {
		log.Warn("user id not found in the ctx")
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Warn("failed to convert user id to int", "error", err)
		return nil, status.Error(codes.Internal, "Something went wrong")
	}

	return &userID, nil
}
