package domain

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}
