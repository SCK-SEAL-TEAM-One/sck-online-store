package auth_test

import (
	"context"
	"store-service/internal/auth"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockJWTTokenManager struct {
	mock.Mock
}

func (manager *mockJWTTokenManager) Generate(user auth.UserPayload, ttl time.Duration) (string, error) {
	args := manager.Called(user, ttl)
	return args.String(0), args.Error(1)
}

func (manager *mockJWTTokenManager) Validate(tokenString string) (auth.Claims, error) {
	args := manager.Called(tokenString)
	return args.Get(0).(auth.Claims), args.Error(1)
}

type mockUserRepository struct {
	mock.Mock
}

func (repo *mockUserRepository) FindByUsername(ctx context.Context, username string) (auth.User, error) {
	args := repo.Called(ctx, username)
	return args.Get(0).(auth.User), args.Error(1)
}

func (repo *mockUserRepository) FindByID(ctx context.Context, uid int) (auth.UserPayload, error) {
	args := repo.Called(ctx, uid)
	return args.Get(0).(auth.UserPayload), args.Error(1)
}

type MockPasswordHelper struct {
	mock.Mock
}

func (helper *MockPasswordHelper) CheckPasswordHash(password, hashed string) bool {
	args := helper.Called(password, hashed)
	return args.Bool(0)
}

func (helper *MockPasswordHelper) HashPassword(password string) (string, error) {
	args := helper.Called(password)
	return args.String(0), args.Error(1)
}
