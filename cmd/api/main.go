package main

import (
	"errors"
	"github.com/AlexShmak/wb_test_task_l0/internal/auth"
	"log"
	"log/slog"

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
		log.Fatalf("failed to load config: %v", err)
	}

	// setup slogLogger
	slogLogger := logger.SetupLogger(cfg.Environment)

	// connect to database
	regularDB, err := db.Connect(cfg)
	if err != nil {
		log.Panic("Could not connect to database.")
	}

	// setup storage
	slogLogger.Info("Running database migrations...")
	m, err := migrate.New(
		"file://cmd/migrations",
		cfg.Database.URL,
	)
	if err != nil {
		log.Fatalf("Could not create migration instance: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Could not apply migrations: %v", err)
	}
	slogLogger.Info("Migrations applied successfully.")
	postgresStorage := storage.NewPostgresStorage(regularDB)

	// setup router
	jwtService := auth.NewJWTService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)
	r := router.NewRouter(postgresStorage, slogLogger, jwtService)
	if err := r.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatalf("Error starting r: %v", err)
	}

	slogLogger.Info("Router started.", slog.String("env", cfg.Environment))
}
