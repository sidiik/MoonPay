package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sidiik/moonpay/api_gateway/internal/constants"
	"github.com/sidiik/moonpay/api_gateway/internal/grpc_clients/auth/authpb"
	"github.com/sidiik/moonpay/api_gateway/internal/infra/config"
	"github.com/sidiik/moonpay/api_gateway/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Authenticate(authClient authpb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			slog.Warn("Missing Authorization header")
			pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			slog.Warn("Authorization header does not start with Bearer")
			pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
			c.Abort()
			return
		}

		accessToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if accessToken == "" {
			slog.Warn("Empty access token")
			pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
			c.Abort()
			return
		}

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
			return []byte(config.AppConfig.AccessTokenSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			slog.Warn("Failed to parse JWT", "error", err)
			pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
			c.Abort()
			return
		}

		if !token.Valid {
			slog.Warn("Invalid JWT token")
			pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Warn("Failed to cast JWT claims to MapClaims")
			pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
			c.Abort()
			return
		}

		email, ok := claims["sub"].(string)
		if !ok {
			slog.Warn("JWT claim 'sub' is missing or invalid")
			pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
			c.Abort()
			return
		}

		slog.Info("Authenticating user", "email", email)

		resp, err := authClient.GetUserByEmail(ctx, &authpb.GetUserByEmailRequest{
			Email: email,
		})

		if err != nil {
			s, _ := status.FromError(err)
			slog.Error("Failed to get user from auth service", "email", email, "grpc_code", s.Code(), "error", s.Message())
			pkg.SendResponse(c, http.StatusUnauthorized, s.Code().String(), "", nil, errors.New(s.Message()))
			c.Abort()
			return
		}

		slog.Info("User authenticated successfully", "email", email, "user_id", resp.Id)
		c.Set("user", resp)
		c.Next()
	}
}
