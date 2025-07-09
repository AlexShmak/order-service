package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			h.Logger.Error("no access token found in cookies", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token, err := h.JWTService.ValidateAccessToken(tokenString)
		if err != nil || !token.Valid {
			h.Logger.Error("invalid token", slog.String("error", err.Error()))
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
				return
			}
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
