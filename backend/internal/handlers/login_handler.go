package handlers

import (
	"github.com/AlexShmak/order-service/internal/storage"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) LoginHandler(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginRequest); err != nil {
		h.Logger.Error("cannot bind JSON", slog.String("error", err.Error()))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	user, err := h.Storage.Users.GetByEmail(c.Request.Context(), loginRequest.Email)
	if err != nil {
		h.Logger.Error("user not found", slog.String("error", err.Error()))
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if !user.CheckPasswordHash(loginRequest.Password) {
		h.Logger.Error("invalid password")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	accessTokenString, refreshTokenString, err := h.JWTService.GenerateTokens(user.ID)
	if err != nil {
		h.Logger.Error("failed to generate tokens", slog.String("error", err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	refreshToken := &storage.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := h.Storage.Tokens.Create(c.Request.Context(), refreshToken); err != nil {
		h.Logger.Error("failed to save refresh token", slog.String("error", err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.SetCookie("access_token", accessTokenString, 15*60, "/", "", false, true)
	c.SetCookie("refresh_token", refreshTokenString, 7*24*60*60, "/", "", false, true)

	h.Logger.Info("user logged in", slog.Int64("id", user.ID))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "logged in"})
}
