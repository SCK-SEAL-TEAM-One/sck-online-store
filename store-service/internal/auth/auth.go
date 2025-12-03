package auth

import (
	"errors"
	"store-service/internal/user"
	"time"
)

type AuthInterface interface {
	Login(username, password string) (TokenPair, error)
	GetAccessToken(claim Claims, ttl time.Duration) (string, error)
	GetRefreshToken(claim Claims, ttl time.Duration) (string, error)
	ValidateToken(token string) (Claims, error)
}

type AuthService struct {
	UserRepository  UserRepository
	JWTTokenManager JWTTokenManagerInterface
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

func (service AuthService) Login(username, password string) (TokenPair, error) {
	userInfo, err := service.UserRepository.FindByUsername(username)
	if err != nil {
		return TokenPair{}, ErrUserNotFound
	}

	if !user.CheckPasswordHash(password, userInfo.Password) {
		return TokenPair{}, ErrInvalidCredentials
	}

	accessTokenTtl := time.Hour            // 1 hour
	refreshTokenTtl := 24 * time.Hour * 30 // 30 days

	claims := Claims{
		UserID:    userInfo.ID,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Username:  userInfo.Username,
	}

	accessToken, err := service.GetAccessToken(claims, accessTokenTtl)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := service.GetRefreshToken(claims, refreshTokenTtl)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (service AuthService) GetAccessToken(claim Claims, ttl time.Duration) (string, error) {
	return service.JWTTokenManager.Generate(claim.ToPayload(), ttl)
}

func (service AuthService) GetRefreshToken(claim Claims, ttl time.Duration) (string, error) {
	return service.JWTTokenManager.Generate(claim.ToPayload(), ttl)
}

func (service AuthService) ValidateToken(token string) (Claims, error) {
	return service.JWTTokenManager.Validate(token)
}
