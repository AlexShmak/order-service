package cache

import (
	"context"
	"github.com/AlexShmak/order-service/internal/storage"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	Orders interface {
		Get(context.Context, string) (*storage.Order, error)
		Set(context.Context, *storage.Order) error
	}
}

func NewRedisStorage(rdb *redis.Client) *RedisStorage {
	return &RedisStorage{
		Orders: &RedisOrders{rdb: rdb},
	}
}
