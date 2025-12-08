package middleware

import (
	"errors"
	"net/http"
	"store-service/internal/auth"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoAuthHeader = errors.New("No authorization header provided.")
	ErrInvalidAuth  = errors.New("Invalid authorization header format.")
	ErrInvalidToken = errors.New("Invalid or expired token.")
)

func ParseAndValidateAuthToken(signature, authHeader string) (*auth.Claims, error) {
	claims := &auth.Claims{}

	if authHeader == "" {
		return nil, ErrNoAuthHeader
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, ErrInvalidAuth
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(signature), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func AuthUser(signature string) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.Request.Header.Get("Authorization")
		claims, err := ParseAndValidateAuthToken(signature, authHeader)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": err.Error(),
			})
			context.Abort()
			return
		}

		context.Set("userID", claims.UserID)

		context.Next()
	}
}
