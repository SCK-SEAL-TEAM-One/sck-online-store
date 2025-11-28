package auth

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    int    `json:"user_id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Username  string `json:"username" db:"username"`
	jwt.RegisteredClaims
}
