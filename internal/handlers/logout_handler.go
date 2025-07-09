package handlers

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) LogoutHandler(c *gin.Context) {
	refreshTokenString, err := c.Cookie("refresh_token")
	if err == nil {
		if err := h.Storage.Tokens.Delete(c.Request.Context(), refreshTokenString); err != nil {
			h.Logger.Error("failed to delete refresh token on logout", slog.String("error", err.Error()))
		}
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	h.Logger.Info("user logged out")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "logged out"})
}
