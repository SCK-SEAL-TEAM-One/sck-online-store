package auth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTTokenManagerInterface interface {
	Generate(user UserPayload, ttl time.Duration) (string, error)
	Validate(tokenString string) (Claims, error)
}

type JWTTokenManager struct {
	SecretKey string
}

func (tokenManager *JWTTokenManager) Generate(user UserPayload, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenManager.SecretKey))
}

func (tokenManager *JWTTokenManager) Validate(tokenString string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenManager.SecretKey), nil
	})

	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return *claims, nil
}
