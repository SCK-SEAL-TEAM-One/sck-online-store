package middleware_test

import (
	"store-service/internal/auth"
	"store-service/internal/middleware"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func Test_ParseAndValidateAuthToken_Should_Return_Error_NoHeader_If_No_Bearer_Token(t *testing.T) {
	// Arrange
	expected := middleware.ErrNoAuthHeader
	signature := "secret"
	authHeader := ""

	// Act
	actual, err := middleware.ParseAndValidateAuthToken(signature, authHeader)

	// Assert
	assert.Nil(t, actual)
	assert.Equal(t, expected, err)
}

func Test_ParseAndValidateAuthToken_Should_Return_Error_InvalidAuth_If_Header_Is_Not_Bearer(t *testing.T) {
	// Arrange
	expected := middleware.ErrInvalidAuth
	signature := "secret"
	authHeader := "NotBearer token"

	// Act
	actual, err := middleware.ParseAndValidateAuthToken(signature, authHeader)

	// Assert
	assert.Nil(t, actual)
	assert.Equal(t, expected, err)
}

func Test_ParseAndValidateAuthToken_Should_Return_Error_InvalidToken_If_Token_Is_Expired(t *testing.T) {
	// Arrange
	userID := 1
	firstName := "Nattapon"
	lastName := "Srisombat"
	username := "nattapon.s"
	expected := middleware.ErrInvalidToken
	signature := "my-secret-key"

	expiredClaims := &auth.Claims{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte(signature))
	authHeader := "Bearer " + expiredTokenString

	// Act
	actual, err := middleware.ParseAndValidateAuthToken(signature, authHeader)

	// Assert
	assert.Nil(t, actual)
	assert.Equal(t, expected, err)
}

func Test_ParseAndValidateAuthToken_Should_Return_Error_InvalidToken_If_Sign_With_Wrong_Signature(t *testing.T) {
	// Arrange
	expected := middleware.ErrInvalidToken
	signature := "my-secret-key"
	wrongSignature := "wrong-secret-key"
	userID := 2
	firstName := "Jakkrit"
	lastName := "Saengthong"
	username := "jakkrit.s"

	claims := &auth.Claims{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(wrongSignature))
	authHeader := "Bearer " + tokenString

	// Act
	actual, err := middleware.ParseAndValidateAuthToken(signature, authHeader)

	// Assert
	assert.Nil(t, actual)
	assert.Equal(t, expected, err)
}

func Test_ParseAndValidateAuthToken_Should_Return_Error_InvalidToken_If_Sign_With_Wrong_SigningMethod(t *testing.T) {
	// Arrange
	expected := middleware.ErrInvalidToken
	signature := "my-secret-key"
	userID := 3
	firstName := "Pimchanok"
	lastName := "Techawong"
	username := "pimchanok.t"

	claims := &auth.Claims{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims) // Sign with RSA256
	tokenString, _ := token.SignedString([]byte(signature))
	authHeader := "Bearer " + tokenString

	// Act
	actual, err := middleware.ParseAndValidateAuthToken(signature, authHeader)

	// Assert
	assert.Nil(t, actual)
	assert.Equal(t, expected, err)
}

func Test_ParseAndValidateAuthToken_Should_Return_Claims_If_Token_Is_Valid(t *testing.T) {
	// Arrange
	userID := 1
	firstName := "Nattapon"
	lastName := "Srisombat"
	username := "nattapon.s"
	signature := "my-secret-key"
	expected := userID

	claimsInput := &auth.Claims{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsInput)
	tokenString, _ := token.SignedString([]byte(signature))
	authHeader := "Bearer " + tokenString

	// Act
	actual, err := middleware.ParseAndValidateAuthToken(signature, authHeader)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual.UserID)
}
