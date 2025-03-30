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
	claims := jwt.MapClaims{
		"userId":    userId,
		"issuer":    "go-jwt-auth-service",
		"expiresAt": jwt.NewNumericDate(time.Now().Add(expired)),
		"issuedAt":  jwt.NewNumericDate(time.Now()),
		"subject":   "auth_token",
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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode signing yang digunakan adalah RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var tokenClaim domain.TokenClaims

		if userId, ok := claims["userId"].(float64); ok {
			tokenClaim.UserID = uint(userId)
		} else {
			return nil, fmt.Errorf("claim 'userId' invalid")
		}

		return &tokenClaim, nil
	}

	return nil, fmt.Errorf("invalid token")
}
