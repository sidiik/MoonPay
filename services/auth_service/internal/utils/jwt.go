package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sidiik/moonpay/auth_service/internal/infra/config"
)

func GenerateAccessToken(email string) (*string, error) {
	tokenStr, err := GenerateToken(email, false)
	return tokenStr, err
}

func GenerateRefreshToken(email string) (*string, error) {
	tokenStr, err := GenerateToken(email, true)
	return tokenStr, err
}

func GenerateToken(email string, isRefresh bool) (*string, error) {

	var (
		jwtExpire int64
		jwtSecret []byte
		err       error
	)

	switch isRefresh {
	case true:
		jwtExpire, err = strconv.ParseInt(config.AppConfig.RefreshTokenExpire, 10, 32)
		if err != nil {
			return nil, err
		}

		jwtSecret = []byte(config.AppConfig.RefreshTokenSecret)
	default:
		jwtExpire, err = strconv.ParseInt(config.AppConfig.AccessTokenExpire, 10, 32)
		if err != nil {
			return nil, err
		}

		jwtSecret = []byte(config.AppConfig.AccessTokenSecret)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * time.Duration(jwtExpire)).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)

	if err != nil {
		return nil, err
	}

	return &tokenStr, nil
}
