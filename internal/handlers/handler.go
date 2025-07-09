package handlers

import (
	"github.com/AlexShmak/wb_test_task_l0/internal/auth"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"log/slog"
)

type Handler struct {
	Storage    *storage.PostgresStorage
	Logger     *slog.Logger
	JWTService *auth.JWTService
}

func NewHandler(s *storage.PostgresStorage, l *slog.Logger, jwt *auth.JWTService) *Handler {
	return &Handler{
		Storage:    s,
		Logger:     l,
		JWTService: jwt,
	}
}
