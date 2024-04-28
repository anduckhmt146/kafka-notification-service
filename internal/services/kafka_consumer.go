package services

import (
	"fmt"

	"github.com/IBM/sarama"
)

type IMessageConsumer interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

type MessageConsumer struct {
	Ready chan bool
}

func NewMessageConsumer() *MessageConsumer {
	return &MessageConsumer{
		Ready: make(chan bool),
	}
}

func (mc *MessageConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(mc.Ready)
	return nil
}

func (mc *MessageConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (mc *MessageConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Message claimed: value = %s, timestamp = %v, topic = %s\n", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}
