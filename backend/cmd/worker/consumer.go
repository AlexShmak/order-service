package worker

import (
	"fmt"
	"github.com/AlexShmak/wb_test_task_l0/internal/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func StartWorker(cfg *config.Config) {
	msgCnt := 0

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerConn, err := sarama.NewConsumer(cfg.Kafka.Brokers, consumerConfig)

	if err != nil {
		panic(err)
	}
	defer func(consumerConn sarama.Consumer) {
		err := consumerConn.Close()
		if err != nil {
			fmt.Printf("Failed to close consumer connection: %v\n", err)
		} else {
			fmt.Println("Consumer connection closed successfully")
		}
	}(consumerConn)

	partitionConsumer, err := consumerConn.ConsumePartition(cfg.Kafka.Topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	defer func(partitionConsumer sarama.PartitionConsumer) {
		err := partitionConsumer.Close()
		if err != nil {
			fmt.Printf("Failed to close partition consumer: %v\n", err)
		} else {
			fmt.Println("Partition consumer closed successfully")
		}
	}(partitionConsumer)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-partitionConsumer.Errors():
				fmt.Println("Consumer error:", err)
			case msg := <-partitionConsumer.Messages():
				msgCnt++
				fmt.Printf("Received order: %s (count=%d)\n", string(msg.Value), msgCnt)
			case <-sigCh:
				fmt.Println("Shutting down consumer...")
				doneCh <- struct{}{}
				return
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCnt, "messages. Exiting.")
}
