package main

import (
	"errors"
	"github.com/AlexShmak/wb_test_task_l0/cmd/worker"
	"github.com/AlexShmak/wb_test_task_l0/internal/auth"
	"github.com/AlexShmak/wb_test_task_l0/internal/kafka"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage/cache"
	"log/slog"
	"os"

	"github.com/AlexShmak/wb_test_task_l0/internal/config"
	"github.com/AlexShmak/wb_test_task_l0/internal/db"
	"github.com/AlexShmak/wb_test_task_l0/internal/logger"
	"github.com/AlexShmak/wb_test_task_l0/internal/router"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// read config
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// setup logger
	slogLogger := logger.SetupLogger(cfg.Environment)

	// connect to database
	regularDB, err := db.Connect(cfg)
	if err != nil {
		slogLogger.Error("Could not connect to database", "error", err)
		os.Exit(1)
	}

	// setup storage
	slogLogger.Info("Running database migrations...")
	m, err := migrate.New(
		"file://cmd/migrations",
		cfg.Database.URL,
	)
	if err != nil {
		slogLogger.Error("Could not create migration instance", "error", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slogLogger.Error("Could not apply migrations", "error", err)
		os.Exit(1)
	}
	slogLogger.Info("Migrations applied successfully.")
	postgresStorage := storage.NewPostgresStorage(regularDB)

	// redisCache setup
	redisClient := cache.NewRedisClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	redisCache := cache.NewRedisStorage(redisClient)

	// start worker
	go worker.StartWorker(cfg, postgresStorage, slogLogger)

	// setup kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg, slogLogger)
	if err != nil {
		slogLogger.Error("failed to create kafka producer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			slogLogger.Error("failed to close kafka producer", "error", err)
		}
	}()

	// setup router
	jwtService := auth.NewJWTService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)
	r := router.NewRouter(postgresStorage, slogLogger, jwtService, cfg, kafkaProducer, redisCache)
	if err := r.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		slogLogger.Error("Error starting r", "error", err)
		os.Exit(1)
	}

	slogLogger.Info("Router started.", slog.String("env", cfg.Environment))

}
