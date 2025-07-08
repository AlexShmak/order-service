package storage

import (
	"database/sql"
	"errors"
	"github.com/AlexShmak/wb_test_task_l0/internal/models"
)

type UsersRepository struct {
	db *sql.DB
}

func (r *UsersRepository) Create(user *models.User) error {
	query := `INSERT INTO orders_service.users (id, name, email, password) VALUES ($1, $2, $3, $4)`
	if _, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.PasswordHash); err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) Delete(user *models.User) error {
	query := `DELETE FROM orders_service.users WHERE id = $1`
	if _, err := r.db.Exec(query, user.ID); err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password FROM orders_service.users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
