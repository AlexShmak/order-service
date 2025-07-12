package kafka

import (
	"github.com/AlexShmak/wb_test_task_l0/internal/config"
	"github.com/IBM/sarama"
	"log/slog"
)

type Producer struct {
	SyncProducer sarama.SyncProducer
	logger       *slog.Logger
}

func NewProducer(cfg *config.Config, logger *slog.Logger) (*Producer, error) {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5

	syncProducer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, producerConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{SyncProducer: syncProducer, logger: logger}, nil
}

func (p *Producer) PushOrderToQueue(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := p.SyncProducer.SendMessage(msg)
	if err != nil {
		return err
	}
	p.logger.Info("Message sent", "partition", partition, "offset", offset)
	return nil
}

func (p *Producer) Close() error {
	if err := p.SyncProducer.Close(); err != nil {
		p.logger.Error("Failed to close producer", "error", err)
		return err
	}
	p.logger.Info("Producer closed successfully")
	return nil
}
