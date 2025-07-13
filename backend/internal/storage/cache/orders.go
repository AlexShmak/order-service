package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlexShmak/order-service/internal/storage"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisOrders struct {
	rdb *redis.Client
}

func (r *RedisOrders) Set(ctx context.Context, order *storage.Order) error {
	cacheKey := fmt.Sprintf("order-%v", order.OrderUID)
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.rdb.SetEx(ctx, cacheKey, data, time.Hour).Err()
}

func (r *RedisOrders) Get(ctx context.Context, uid string) (*storage.Order, error) {
	cacheKey := fmt.Sprintf("order-%v", uid)
	data, err := r.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var order storage.Order
	if data != "" {
		if err := json.Unmarshal([]byte(data), &order); err != nil {
			return nil, err
		}
	}
	return &order, nil
}
