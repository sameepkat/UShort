package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID uint
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
