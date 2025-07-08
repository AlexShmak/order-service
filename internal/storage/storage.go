package storage

import (
	"database/sql"

	"github.com/AlexShmak/wb_test_task_l0/internal/models"
)

type PostgresStorage struct {
	Users interface {
		Create(user *models.User) error
		Delete(user *models.User) error
		GetByID(id int64) (*models.User, error)
	}
	OrdersAdmin interface {
		Create(order *models.Order) error
		DeleteByID(id int64) error
		GetByID(id int64) (*models.Order, error)
	}
	OrdersRegular interface {
		GetByID(id int64) (*models.Order, error)
	}
}

func NewPostgresStorage(adminDB *sql.DB, regularDB *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		Users:         &UsersRepository{db: adminDB},
		OrdersAdmin:   &OrdersRepositoryAdmin{db: adminDB},
		OrdersRegular: &OrdersRepositoryRegular{db: regularDB},
	}
}
