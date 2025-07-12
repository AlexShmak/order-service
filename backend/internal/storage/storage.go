package storage

import (
	"context"
	"database/sql"
)

type PostgresStorage struct {
	Users interface {
		Create(context.Context, *User) error
		GetByEmail(context.Context, string) (*User, error)
	}
	Orders interface {
		GetByID(context.Context, int64, int64) (*Order, error)
		Create(ctx context.Context, order *Order) error
	}
	Tokens interface {
		Create(context.Context, *RefreshToken) error
		Delete(context.Context, string) error
		GetByToken(context.Context, string) (*RefreshToken, error)
	}
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		Users:  &UsersRepository{db: db},
		Orders: &OrdersRepository{db: db},
		Tokens: &TokensRepository{db: db},
	}
}
