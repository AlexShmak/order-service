package kafka

import (
	"github.com/AlexShmak/wb_test_task_l0/internal/config"
	"github.com/IBM/sarama"
	"log"
)

func PushOrderToQueue(topic string, message []byte, cfg *config.Config) error {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, producerConfig)

	if err != nil {
		return err
	}
	defer func(producer sarama.SyncProducer) {
		err := producer.Close()
		if err != nil {
			log.Printf("Failed to close producer: %v", err)
		} else {
			log.Println("Producer closed successfully")
		}
	}(producer)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	return nil
}
