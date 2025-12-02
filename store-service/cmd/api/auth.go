package api

import (
	"encoding/json"
	"log"
	"net/http"
	"store-service/internal/auth"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthAPI struct {
	AuthService auth.AuthInterface
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (api AuthAPI) RefreshTokenHandler(context *gin.Context) {
	cookie, err := context.Cookie("refreshToken")
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   "no_token",
			"message": "No refresh token found. Please log in again.",
		})
	}

	claims, err := api.AuthService.ValidateToken(cookie)
	if err != nil {
		context.SetCookie("refreshToken", "", -1, "/", "", false, true)

		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_refresh_token",
			"message": "Refresh token is expired or invalid. Please log in again.",
		})
		return
	}

	accessTokenTtl := time.Hour            // 1 hour
	refreshTokenTtl := 24 * time.Hour * 30 // 30 days

	accessToken, err := api.AuthService.GetAccessToken(claims, accessTokenTtl)
	if err != nil {
		context.SetCookie("refreshToken", "", -1, "/", "", false, true)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
	}
	refreshToken, err := api.AuthService.GetRefreshToken(claims, refreshTokenTtl)
	if err != nil {
		context.SetCookie("refreshToken", "", -1, "/", "", false, true)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
	}

	context.SetCookie("refreshToken", refreshToken, 24*3600*30, "/", "", false, true)
	context.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"message":      "Access token refreshed successfully.",
	})
}

func (api AuthAPI) LoginHandler(context *gin.Context) {
	var req LoginRequest
	decoder := json.NewDecoder(context.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid input.",
		})
		log.Printf("bad request %s", err.Error())
		return
	}

	if req.Username == "" || req.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Both username and password are required.",
		})
		return
	}

	tokens, err := api.AuthService.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials, auth.ErrUserNotFound:
			context.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_credentials",
				"message": "Invalid email or password.",
			})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate token",
			})
		}
		return
	}

	context.SetCookie("refreshToken", tokens.RefreshToken, 24*3600*30, "/", "", false, true)
	context.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
		"message":      "Access token refreshed successfully.",
	})
}
