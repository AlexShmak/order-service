package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/AlexShmak/order-service/internal/config"
)

func Connect(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.GetPostgresDSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.Database.MaxOpenConns)
	db.SetMaxIdleConns(config.Database.MaxIdleConns)
	db.SetConnMaxLifetime(config.Database.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.Database.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
