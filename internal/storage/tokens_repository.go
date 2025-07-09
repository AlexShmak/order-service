package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type RefreshToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
type TokensRepository struct {
	db *sql.DB
}

func (s *TokensRepository) Create(ctx context.Context, token *RefreshToken) error {
	query := `
		INSERT INTO orders_service.refresh_tokens (user_id, token, expires_at) 
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	if err := s.db.QueryRowContext(ctx, query, token.UserID, token.Token, token.ExpiresAt).Scan(&token.ID, &token.CreatedAt); err != nil {
		return fmt.Errorf("could not create refresh token: %w", err)
	}
	return nil
}

func (s *TokensRepository) GetByToken(ctx context.Context, tokenString string) (*RefreshToken, error) {
	var token RefreshToken

	query := `
		SELECT * FROM orders_service.refresh_tokens WHERE token = $1 AND expires_at > NOW()
	`
	if err := s.db.QueryRowContext(ctx, query, tokenString).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, fmt.Errorf("could not get refresh token by token string: %w", err)
	}
	return &token, nil
}

func (s *TokensRepository) Delete(ctx context.Context, tokenString string) error {
	query := `
		DELETE FROM refresh_tokens WHERE token = $1
	`
	if _, err := s.db.ExecContext(ctx, query, tokenString); err != nil {
		return fmt.Errorf("could not delete refresh token: %w", err)
	}
	return nil
}
