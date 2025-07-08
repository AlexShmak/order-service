package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/AlexShmak/wb_test_task_l0/internal/config"
)

func Connect(config *config.Config) (*sql.DB, *sql.DB, error) {
	regularDB, err := sql.Open("postgres", config.GetRegularPostgresDSN())
	if err != nil {
		return nil, nil, err
	}
	adminDB, err := sql.Open("postgres", config.GetAdminPostgresDSN())
	if err != nil {
		return nil, nil, err
	}

	regularDB.SetMaxOpenConns(config.Database.MaxOpenConns)
	regularDB.SetMaxIdleConns(config.Database.MaxIdleConns)
	regularDB.SetConnMaxLifetime(config.Database.ConnMaxLifetime)
	regularDB.SetConnMaxIdleTime(config.Database.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = regularDB.PingContext(ctx); err != nil {
		return nil, nil, err
	}
	if err = adminDB.PingContext(ctx); err != nil {
		return nil, nil, err
	}

	return regularDB, adminDB, nil
}
