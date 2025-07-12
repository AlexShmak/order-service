package handlers

import (
	"github.com/AlexShmak/wb_test_task_l0/internal/auth"
	"github.com/AlexShmak/wb_test_task_l0/internal/config"
	"github.com/AlexShmak/wb_test_task_l0/internal/kafka"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage/cache"
	"log/slog"
)

type Handler struct {
	Config        *config.Config
	Storage       *storage.PostgresStorage
	Logger        *slog.Logger
	JWTService    *auth.JWTService
	KafkaProducer *kafka.Producer
	Cache         *cache.RedisStorage
}

func NewHandler(
	storage *storage.PostgresStorage,
	logger *slog.Logger,
	jwtService *auth.JWTService,
	cfg *config.Config,
	producer *kafka.Producer,
	redisCache *cache.RedisStorage,
) *Handler {
	return &Handler{
		Config:        cfg,
		Storage:       storage,
		Logger:        logger,
		JWTService:    jwtService,
		KafkaProducer: producer,
		Cache:         redisCache,
	}
}
