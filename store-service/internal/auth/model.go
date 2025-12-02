package auth

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        int    `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Password  string `json:"password" db:"password"`
}

type UserPayload struct {
	UserID    int    `json:"user_id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Username  string `json:"username" db:"username"`
}

type Claims struct {
	UserID    int    `json:"user_id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Username  string `json:"username" db:"username"`
	jwt.RegisteredClaims
}

func (claim Claims) ToPayload() UserPayload {
	return UserPayload{
		UserID:    claim.UserID,
		FirstName: claim.FirstName,
		LastName:  claim.LastName,
		Username:  claim.Username,
	}
}
