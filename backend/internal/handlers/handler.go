package handlers

import (
	"github.com/AlexShmak/wb_test_task_l0/internal/auth"
	"github.com/AlexShmak/wb_test_task_l0/internal/config"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"log/slog"
)

type Handler struct {
	Config     *config.Config
	Storage    *storage.PostgresStorage
	Logger     *slog.Logger
	JWTService *auth.JWTService
}

func NewHandler(s *storage.PostgresStorage, l *slog.Logger, jwt *auth.JWTService, cfg *config.Config) *Handler {
	return &Handler{
		Config:     cfg,
		Storage:    s,
		Logger:     l,
		JWTService: jwt,
	}
}
