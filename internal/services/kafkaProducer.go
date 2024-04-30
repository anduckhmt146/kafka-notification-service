package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
}

func NewKafkaProducer(kafkaBrokers string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	config.Version = sarama.V2_5_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100

	sarama.Logger = log.New(os.Stdout, "[sarama-producer] ", log.LstdFlags)

	// Create the producer
	producer, err := sarama.NewSyncProducer(strings.Split(kafkaBrokers, ","), config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}

	return &KafkaProducer{Producer: producer}, nil
}

func (kp *KafkaProducer) Cleanup() error {
	if err := kp.Producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %w", err)
	}
	return nil
}
