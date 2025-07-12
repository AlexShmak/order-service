package handlers

import (
	"database/sql"
	"errors"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) handleTokenRefresh(c *gin.Context, originalUserID int64) (int64, bool) {
	oldRefreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		h.Logger.Error("no refresh token found in cookies", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
		return 0, false
	}

	if _, err = h.Storage.Tokens.GetByToken(c.Request.Context(), oldRefreshTokenString); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.Warn("refresh token not found in db", slog.String("token", oldRefreshTokenString))
		} else {
			h.Logger.Error("failed to validate refresh token", slog.String("error", err.Error()))
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return 0, false
	}

	token, _, _ := new(jwt.Parser).ParseUnverified(oldRefreshTokenString, jwt.MapClaims{})
	refreshUserID, err := h.JWTService.GetUserIdFromToken(token)
	if err != nil {
		h.Logger.Error("could not get user ID from old refresh token", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return 0, false
	}

	if refreshUserID != originalUserID {
		h.Logger.Error("token refresh user ID mismatch", slog.Int64("original_user", originalUserID), slog.Int64("refresh_user", refreshUserID))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token mismatch"})
		return 0, false
	}

	if err := h.Storage.Tokens.Delete(c.Request.Context(), oldRefreshTokenString); err != nil {
		h.Logger.Error("failed to delete old refresh token", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not refresh session"})
		return 0, false
	}

	newAccessTokenString, newRefreshTokenString, err := h.JWTService.GenerateTokens(refreshUserID)
	if err != nil {
		h.Logger.Error("failed to generate new tokens", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not refresh session"})
		return 0, false
	}

	newRefreshToken := &storage.RefreshToken{
		UserID:    refreshUserID,
		Token:     newRefreshTokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := h.Storage.Tokens.Create(c.Request.Context(), newRefreshToken); err != nil {
		h.Logger.Error("failed to save new refresh token", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not refresh session"})
		return 0, false
	}

	c.SetCookie("access_token", newAccessTokenString, 15*60, "/", "", false, true)
	c.SetCookie("refresh_token", newRefreshTokenString, 7*24*60*60, "/", "", false, true)
	h.Logger.Info("tokens refreshed successfully", slog.Int64("userID", refreshUserID))

	return refreshUserID, true
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			h.Logger.Info("access token not found, attempting refresh with refresh token")

			refreshTokenString, refreshErr := c.Cookie("refresh_token")
			if refreshErr != nil {
				h.Logger.Error("no access or refresh token found in cookies", slog.String("error", refreshErr.Error()))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			token, _, parseErr := new(jwt.Parser).ParseUnverified(refreshTokenString, jwt.MapClaims{})
			if parseErr != nil {
				h.Logger.Error("could not parse refresh token", slog.String("error", parseErr.Error()))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}

			userID, idErr := h.JWTService.GetUserIdFromToken(token)
			if idErr != nil {
				h.Logger.Error("could not get user ID from refresh token", slog.String("error", idErr.Error()))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
				return
			}

			if refreshedUserID, ok := h.handleTokenRefresh(c, userID); ok {
				c.Set("userId", refreshedUserID)
				c.Next()
			}
			return
		}

		token, err := h.JWTService.ValidateAccessToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				h.Logger.Info("access token expired, attempting refresh")

				userID, idErr := h.JWTService.GetUserIdFromToken(token)
				if idErr != nil {
					h.Logger.Error("could not get user ID from expired token", slog.String("error", idErr.Error()))
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
					return
				}

				if refreshedUserID, ok := h.handleTokenRefresh(c, userID); ok {
					c.Set("userId", refreshedUserID)
					c.Next()
				}
				return
			}

			h.Logger.Error("invalid token", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userId, err := h.JWTService.GetUserIdFromToken(token)
		if err != nil {
			h.Logger.Error("could not get user ID from token", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
