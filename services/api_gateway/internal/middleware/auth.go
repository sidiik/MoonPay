package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sidiik/moonpay/api_gateway/internal/constants"
	"github.com/sidiik/moonpay/api_gateway/internal/infra/config"
	"github.com/sidiik/moonpay/api_gateway/pkg"
	"google.golang.org/grpc/codes"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		userID, ok := claims["sub"].(string)
		if !ok {
			slog.Warn("JWT claim 'sub' is missing or invalid", "user-id", claims["sub"])
			pkg.SendResponse(c, http.StatusUnauthorized, codes.Unauthenticated.String(), constants.ErrUnauthorized, nil, nil)
			c.Abort()
			return
		}

		slog.Info("User authenticated successfully", "userID", userID)
		c.Set("userID", userID)
		c.Next()
	}
}
