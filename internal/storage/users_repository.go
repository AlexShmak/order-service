package storage

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:""`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	passwordHash []byte
}

func (u *User) hashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.passwordHash = hash

	return nil
}

func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.passwordHash, []byte(password))
	return err == nil
}

type UsersRepository struct {
	db *sql.DB
}

func (s *UsersRepository) Create(ctx context.Context, user *User) error {

	if err := user.hashPassword(); err != nil {
		return err
	}

	query := `
		INSERT INTO orders_service.users (username, password, email)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	err := s.db.QueryRowContext(ctx, query, user.Username, user.passwordHash, user.Email).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT * FROM orders_service.users
		WHERE email = $1
	`
	user := &User{}
	var passwordHash []byte
	if err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &passwordHash, &user.Email, &user.CreatedAt); err != nil {
		return nil, err
	}
	user.passwordHash = passwordHash
	return user, nil
}
