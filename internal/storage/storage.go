package storage

import (
	"database/sql"

	"github.com/AlexShmak/wb_test_task_l0/internal/models"
)

type PostgresStorage struct {
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
		OrdersAdmin:   &OrdersRepositoryAdmin{db: adminDB},
		OrdersRegular: &OrdersRepositoryRegular{db: regularDB},
	}
}
