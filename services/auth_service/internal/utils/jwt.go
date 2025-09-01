package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sidiik/moonpay/auth_service/internal/infra/config"
)

func GenerateAccessToken(id uint) (*string, error) {
	tokenStr, err := GenerateToken(id, false)
	return tokenStr, err
}

func GenerateRefreshToken(id uint) (*string, error) {
	tokenStr, err := GenerateToken(id, true)
	return tokenStr, err
}

func GenerateToken(id uint, isRefresh bool) (*string, error) {

	var (
		jwtExpire int64
		jwtSecret []byte
		err       error
	)

	idStr := strconv.Itoa(int(id))

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
		"sub": idStr,
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * time.Duration(jwtExpire)).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)

	if err != nil {
		return nil, err
	}

	return &tokenStr, nil
}
