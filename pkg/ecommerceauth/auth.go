package ecommerceauth

import (
	"fmt"
	"math"
	"time"

	"github.com/NATCHAYATP/E-Commerce/config"
	"github.com/NATCHAYATP/E-Commerce/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type ecommerceAuth struct {
	mapClaims *ecommerceMapClaims //payload
	cfg       config.IJwtConfig
}

type ecommerceMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type IEcommerceAuth interface {
	SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *ecommerceAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func NewEcommerceAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IEcommerceAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IEcommerceAuth {
	return &ecommerceAuth{
		cfg: cfg,
		mapClaims: &ecommerceMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecommerce-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IEcommerceAuth {
	return &ecommerceAuth{
		cfg: cfg,
		mapClaims: &ecommerceMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecommerce-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}
