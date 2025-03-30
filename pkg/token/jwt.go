package token

import (
	"category-service/internal/domain"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewJWT(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) Token {
	return &JWT{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (j *JWT) GenerateToken(userId uint, expired time.Duration) (string, error) {
	claims := domain.TokenClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expired)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-jwt-auth-service",
			Subject:   "auth_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWT) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	// Parse token dan verifikasi menggunakan publicKey
	token, err := jwt.ParseWithClaims(tokenString, &domain.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode signing yang digunakan adalah RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*domain.TokenClaims); ok && token.Valid {
		fmt.Println(claims.UserID)
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
