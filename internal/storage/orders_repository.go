package storage

import (
	"database/sql"

	"github.com/AlexShmak/wb_test_task_l0/internal/models"
)

type OrdersRepositoryRegular struct {
	db *sql.DB
}

type OrdersRepositoryAdmin struct {
	db *sql.DB
}

func (r *OrdersRepositoryRegular) GetByID(id int64) (*models.Order, error) {
	var order models.Order
	query := `SELECT * FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Delivery,
		&order.Payment,
		&order.Items,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrdersRepositoryAdmin) GetByID(id int64) (*models.Order, error) {
	var order models.Order
	query := `SELECT * FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Delivery,
		&order.Payment,
		&order.Items,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrdersRepositoryAdmin) Create(order *models.Order) error {
	query := `INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	if _, err := r.db.Exec(query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Delivery,
		order.Payment,
		order.Items,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	); err != nil {
		return err
	}
	return nil
}

func (r *OrdersRepositoryAdmin) DeleteByID(id int64) error {
	query := `DELETE FROM orders WHERE id = $1`
	if _, err := r.db.Exec(query, id); err != nil {
		return err
	}
	return nil
}
