package auth_test

import (
	"database/sql"
	"store-service/internal/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Login_Should_Return_Error_UserNotFound_If_Incorrect_Username(t *testing.T) {
	// Arrange
	expectedResult := auth.TokenPair{}
	expectedError := auth.ErrUserNotFound
	username := "nattaponnn"
	password := "Natta@2025"

	mockUserRepository := new(mockUserRepository)
	mockUserRepository.On("FindByUsername", username).Return(auth.User{}, sql.ErrNoRows)

	mockPasswordHelper := new(MockPasswordHelper)
	mockTokenManager := new(mockJWTTokenManager)

	authService := auth.AuthService{
		UserRepository: mockUserRepository,
		PasswordHelper: mockPasswordHelper,
	}

	// Act
	actual, err := authService.Login(username, password)

	// Assert
	assert.Equal(t, expectedResult, actual)
	assert.Equal(t, expectedError, err)
	mockPasswordHelper.AssertNotCalled(t, "CheckPasswordHash", mock.Anything, mock.Anything)
	mockTokenManager.AssertNotCalled(t, "Generate", mock.Anything, mock.Anything)
}

func Test_Login_Should_Return_Error_InvalidCredentials_If_Incorrect_Password(t *testing.T) {
	// Arrange
	expectedResult := auth.TokenPair{}
	expectedError := auth.ErrInvalidCredentials
	username := "nattapon.s"
	password := "Natta@123" // wrong password

	foundUser := auth.User{
		ID:        1,
		FirstName: "Nattapon",
		LastName:  "Srisombat",
		Username:  "nattapon.s",
		Password:  "$2a$12$jO/6faXH5oll0iKupXYscuzxiN6Qj6REfGgB18WhGBDun/p7wG0Si",
	}
	mockUserRepository := new(mockUserRepository)
	mockUserRepository.On("FindByUsername", username).Return(foundUser, nil)

	mockPasswordHelper := new(MockPasswordHelper)
	mockPasswordHelper.On("CheckPasswordHash", password, foundUser.Password).Return(false)

	mockTokenManager := new(mockJWTTokenManager)

	authService := auth.AuthService{
		UserRepository: mockUserRepository,
		PasswordHelper: mockPasswordHelper,
	}

	// Act
	actual, err := authService.Login(username, password)

	// Assert
	assert.Equal(t, expectedResult, actual)
	assert.Equal(t, expectedError, err)
	mockTokenManager.AssertNotCalled(t, "Generate", mock.Anything, mock.Anything)
}

func Test_Login_Should_Resturen_TokenPair_If_Username_Password_Is_Correct(t *testing.T) {
	// Arrange
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJmaXJzdF9uYW1lIjoiTmF0dGFwb24iLCJsYXN0X25hbWUiOiJTcmlzb21iYXQiLCJ1c2VybmFtZSI6Im5hdHRhcG9uLnMiLCJzdWIiOiIxIiwiaXNzIjoic3RvcmUtc2VydmljZSIsImF1ZCI6InN0b3JlLWNsaWVudCIsImlhdCI6MTczNTczMjgwMCwibmJmIjoxNzM1NzMyODAwLCJleHAiOjE3MzU3MzY0MDB9.mwwOS-vwo_AvRDVMMInzkNp_QYO3i7irdIR228bz_sk"
	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJmaXJzdF9uYW1lIjoiTmF0dGFwb24iLCJsYXN0X25hbWUiOiJTIiwidXNlcm5hbWUiOiJuYXR0YXBvbi5zIiwic3ViIjoiMSIsImlzcyI6InN0b3JlLXNlcnZpY2UiLCJhdWQiOiJzdG9yZS1jbGllbnQiLCJpYXQiOjE3MzU3MzI4MDAsIm5iZiI6MTczNTczMjgwMCwiZXhwIjoxNzM4MzI0ODAwfQ.0XCBau3Ol_NbjnXXtzWVv-K_E-ocexEh0WNoQ06AcvQ"
	expectedResult := auth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	username := "nattapon.s"
	password := "Natta@2025"

	foundUser := auth.User{
		ID:        1,
		FirstName: "Nattapon",
		LastName:  "Srisombat",
		Username:  "nattapon.s",
		Password:  "$2a$12$jO/6faXH5oll0iKupXYscuzxiN6Qj6REfGgB18WhGBDun/p7wG0Si",
	}
	mockUserRepository := new(mockUserRepository)
	mockUserRepository.On("FindByUsername", username).Return(foundUser, nil)

	mockPasswordHelper := new(MockPasswordHelper)
	mockPasswordHelper.On("CheckPasswordHash", password, foundUser.Password).Return(true)

	claims := auth.Claims{
		UserID:    foundUser.ID,
		FirstName: foundUser.FirstName,
		LastName:  foundUser.LastName,
		Username:  foundUser.Username,
	}
	payload := claims.ToPayload()
	mockTokenManager := new(mockJWTTokenManager)
	mockTokenManager.
		On("Generate", payload, mock.AnythingOfType("time.Duration")).
		Return(accessToken, nil).
		Once()

	mockTokenManager.
		On("Generate", payload, mock.AnythingOfType("time.Duration")).
		Return(refreshToken, nil).
		Once()

	authService := auth.AuthService{
		UserRepository:  mockUserRepository,
		PasswordHelper:  mockPasswordHelper,
		JWTTokenManager: mockTokenManager,
	}

	// Act
	actual, err := authService.Login(username, password)

	// Assert
	assert.Equal(t, expectedResult, actual)
	assert.Equal(t, nil, err)
}
