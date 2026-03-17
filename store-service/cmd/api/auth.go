package api

import (
	"encoding/json"
	"log/slog"
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
	ctx := context.Request.Context()
	ip := context.ClientIP()

	cookie, err := context.Cookie("refreshToken")
	if err != nil {
		slog.WarnContext(ctx, "Token refresh failed: no token",
			"log_type", "audit",
			"actor_id", 0,
			"actor_type", "user",
			"action", "refresh_token",
			"resource_type", "session",
			"resource_id", "",
			"ip_address", ip,
		)
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   "no_token",
			"message": "No refresh token found. Please log in again.",
		})
		return
	}

	claims, err := api.AuthService.ValidateToken(cookie)
	if err != nil {
		slog.WarnContext(ctx, "Token refresh failed: invalid token",
			"log_type", "audit",
			"actor_id", 0,
			"actor_type", "user",
			"action", "refresh_token",
			"resource_type", "session",
			"resource_id", "",
			"ip_address", ip,
		)
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
		slog.ErrorContext(ctx, "AuthService.GetAccessToken internal error",
			"log_type", "error",
			"error_code", "TOKEN_GENERATION_FAILED",
			"error_message", err.Error(),
			"user_id", claims.UserID,
			"ip_address", ip,
		)
		context.SetCookie("refreshToken", "", -1, "/", "", false, true)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
		return
	}
	refreshToken, err := api.AuthService.GetRefreshToken(claims, refreshTokenTtl)
	if err != nil {
		slog.ErrorContext(ctx, "AuthService.GetRefreshToken internal error",
			"log_type", "error",
			"error_code", "TOKEN_GENERATION_FAILED",
			"error_message", err.Error(),
			"user_id", claims.UserID,
			"ip_address", ip,
		)
		context.SetCookie("refreshToken", "", -1, "/", "", false, true)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	slog.InfoContext(ctx, "Token refreshed",
		"log_type", "audit",
		"actor_id", claims.UserID,
		"actor_type", "user",
		"action", "refresh_token",
		"resource_type", "session",
		"resource_id", "",
		"ip_address", ip,
	)

	context.SetCookie("refreshToken", refreshToken, 24*3600*30, "/", "", false, true)
	context.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"message":      "Access token refreshed successfully.",
	})
}

func (api AuthAPI) LoginHandler(context *gin.Context) {
	ctx := context.Request.Context()
	ip := context.ClientIP()

	var req LoginRequest
	decoder := json.NewDecoder(context.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		slog.ErrorContext(ctx, "Login bad request",
			"log_type", "error",
			"error_code", "INVALID_REQUEST",
			"error_message", err.Error(),
			"user_id", 0,
			"ip_address", ip,
		)
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid input.",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		slog.WarnContext(ctx, "Login failed: missing credentials",
			"log_type", "audit",
			"actor_id", 0,
			"actor_type", "user",
			"action", "login",
			"resource_type", "session",
			"resource_id", "",
			"ip_address", ip,
		)
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Both username and password are required.",
		})
		return
	}

	tokens, err := api.AuthService.Login(ctx, req.Username, req.Password)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials, auth.ErrUserNotFound:
			slog.WarnContext(ctx, "Login failed: invalid credentials",
				"log_type", "audit",
				"actor_id", 0,
				"actor_type", "user",
				"action", "login",
				"resource_type", "session",
				"resource_id", req.Username,
				"ip_address", ip,
			)
			context.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_credentials",
				"message": "Invalid email or password.",
			})
		default:
			slog.ErrorContext(ctx, "AuthService.Login internal error",
				"log_type", "error",
				"error_code", "LOGIN_INTERNAL_ERROR",
				"error_message", err.Error(),
				"user_id", 0,
				"ip_address", ip,
			)
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate token",
			})
		}
		return
	}

	slog.InfoContext(ctx, "User logged in",
		"log_type", "audit",
		"actor_id", 0,
		"actor_type", "user",
		"action", "login",
		"resource_type", "session",
		"resource_id", req.Username,
		"ip_address", ip,
	)

	context.SetCookie("refreshToken", tokens.RefreshToken, 24*3600*30, "/", "", false, true)
	context.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
		"message":      "Logged in successfully.",
	})
}
