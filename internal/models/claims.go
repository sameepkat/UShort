package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID uint
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
