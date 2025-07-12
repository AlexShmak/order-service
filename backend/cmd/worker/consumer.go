package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexShmak/wb_test_task_l0/internal/config"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"github.com/IBM/sarama"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Consumer struct {
	ready   chan bool
	Storage *storage.PostgresStorage
	logger  *slog.Logger
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.logger.Info(
			"Message claimed",
			"value", string(message.Value),
			"timestamp", message.Timestamp,
			"topic", message.Topic,
			"partition", message.Partition,
		)

		var order storage.Order
		if err := json.Unmarshal(message.Value, &order); err != nil {
			c.logger.Error("Failed to unmarshal message", "error", err)
			continue
		}

		if err := createOrder(context.Background(), &order, c.Storage, c.logger); err != nil {
			c.logger.Error("Failed to create order", "error", err)
		}

		session.MarkMessage(message, "")
	}
	return nil
}

func createOrder(ctx context.Context, order *storage.Order, storage *storage.PostgresStorage, logger *slog.Logger) error {
	if err := storage.Orders.Create(ctx, order); err != nil {
		return fmt.Errorf("failed to create order in storage: %w", err)
	}
	logger.Info("Order created successfully", "order_uid", order.OrderUID)
	return nil
}

func StartWorker(cfg *config.Config, pgStorage *storage.PostgresStorage, logger *slog.Logger) {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerConfig.Version = sarama.V2_8_1_0 // Specify a specific Kafka version

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, "orders-group", consumerConfig)
	if err != nil {
		logger.Error("Error creating consumer group client", "error", err)
		os.Exit(1)
	}

	consumer := &Consumer{
		ready:   make(chan bool),
		Storage: pgStorage,
		logger:  logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, []string{cfg.Kafka.Topic}, consumer); err != nil {
				logger.Error("Error from consumer", "error", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	logger.Info("Consumer is ready")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info("terminating: context cancelled")
	case <-sigterm:
		logger.Info("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = consumerGroup.Close(); err != nil {
		logger.Error("Error closing client", "error", err)
	}
	logger.Info("Consumer closed.")
}
