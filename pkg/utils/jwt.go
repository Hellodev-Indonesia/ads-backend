package utils

import (
	"time"

	"github.com/alex/ads_backend/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID      uint     `json:"user_id"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email string, roles, permissions []string) (string, error) {
	secret := config.GetEnv("JWT_SECRET", "secret")
	expiration, _ := time.ParseDuration(config.GetEnv("JWT_EXPIRATION", "72h"))

	claims := &JWTClaims{
		UserID:      userID,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateCentrifugoToken(userID string) (string, error) {
	secret := config.GetEnv("CENTRIFUGO_TOKEN_SECRET", "centrifugo_secret")

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	secret := config.GetEnv("JWT_SECRET", "secret")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
