package storage

import (
	"context"
	"database/sql"
)

type Users interface {
	Create(context.Context, *User) error
	GetByEmail(context.Context, string) (*User, error)
}
type Orders interface {
	GetByID(context.Context, string, int64) (*Order, error)
	Create(ctx context.Context, order *Order) error
}
type Tokens interface {
	Create(context.Context, *RefreshToken) error
	Delete(context.Context, string) error
	GetByToken(context.Context, string) (*RefreshToken, error)
}

type PostgresStorage struct {
	Users  Users
	Orders Orders
	Tokens Tokens
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		Users:  &UsersRepository{db: db},
		Orders: &OrdersRepository{db: db},
		Tokens: &TokensRepository{db: db},
	}
}
