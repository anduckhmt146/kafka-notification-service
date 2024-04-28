package services

import (
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/anduckhmt146/kakfa-consumer/internal/constants"
)

type KafkaConsumerGroup struct {
	Client sarama.ConsumerGroup
}

func NewKafkaConsumerGroup(kafkaBrokers, group, version string) KafkaConsumerGroup {
	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = kafkaVersion
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaConfig.ClientID = constants.CLIENT_ID + "_" + group
	kafkaConfig.Consumer.Return.Errors = true
	sarama.Logger = log.New(os.Stdout, "[sarama-logger] ", log.LstdFlags)

	client, err := sarama.NewConsumerGroup(strings.Split(kafkaBrokers, ","), group, kafkaConfig)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	return KafkaConsumerGroup{
		Client: client,
	}
}
