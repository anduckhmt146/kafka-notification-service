package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/anduckhmt146/kakfa-consumer/internal/models"
	"github.com/anduckhmt146/kakfa-consumer/internal/repositories"
	"github.com/spf13/viper"
)

type NotificationService struct {
	Producer sarama.SyncProducer
	UserRepo repositories.IUserRepository
}

type INotificationService interface {
	SendNotification(fromID int, toID int, message string) error
}

func NewNotificationService(producer sarama.SyncProducer, userRepo repositories.IUserRepository) INotificationService {
	return &NotificationService{
		Producer: producer,
		UserRepo: userRepo,
	}
}

func (ns *NotificationService) SendNotification(fromID int, toID int, message string) error {
	fromUser, err := ns.UserRepo.GetUserByID(fromID)
	if err != nil {
		return err
	}

	toUser, err := ns.UserRepo.GetUserByID(toID)
	if err != nil {
		return err
	}

	notification := models.Notification{
		FromID:  fromUser.ID,
		ToID:    toUser.ID,
		Message: message,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: viper.GetString("kafka.topic"),
		Key:   sarama.StringEncoder(strconv.Itoa(toUser.ID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = ns.Producer.SendMessage(msg)
	return err
}
