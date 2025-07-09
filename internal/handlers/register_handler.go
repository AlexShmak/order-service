package handlers

import (
	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) RegisterHandler(c *gin.Context) {
	var user storage.User

	if err := c.BindJSON(&user); err != nil {
		h.Logger.Error("cannot bind JSON", slog.String("error", err.Error()))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.Storage.Users.Create(c.Request.Context(), &user); err != nil {
		h.Logger.Error("cannot create user", slog.String("error", err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	h.Logger.Info("user created", slog.Int64("id", user.ID))
	c.IndentedJSON(http.StatusCreated, gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "created_at": user.CreatedAt})
}
