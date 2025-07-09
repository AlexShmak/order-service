package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Order struct {
	OrderUID          string
	TrackNumber       string
	Entry             string
	Delivery          Delivery
	Payment           Payment
	Items             []Item
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	ShardKey          string
	SmID              int64
	DateCreated       time.Time
	OofShard          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Delivery struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type Payment struct {
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDt    int64
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

type Item struct {
	ChrtID      int
	TrackNumber string
	Price       int
	RID         string
	Name        string
	Sale        int
	Size        string
	TotalPrice  int
	NmID        int
	Brand       string
	Status      int
}

type OrdersRepository struct {
	db *sql.DB
}

func (r *OrdersRepository) GetByID(ctx context.Context, id int64) (*Order, error) {
	var order Order
	var delivery Delivery
	var payment Payment
	var items []Item

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func(tx *sql.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("Rollback error: %v", err)
		}
	}(tx)

	orderQuery := `
        SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
               o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
               d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
               p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt,
               p.bank, p.delivery_cost, p.goods_total, p.custom_fee
        FROM orders_service.orders o
        JOIN orders_service.deliveries d ON o.delivery_data_id = d.id
        JOIN orders_service.payments p ON o.payment_data_id = p.id
        WHERE o.id = $1`

	err = tx.QueryRowContext(ctx, orderQuery, id).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
		&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt,
		&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}
	order.Delivery = delivery
	order.Payment = payment

	itemsQuery := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM orders_service.items WHERE id = $1`
	rows, err := tx.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	order.Items = items

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}
