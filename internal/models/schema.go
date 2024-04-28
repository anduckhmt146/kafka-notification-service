package models

type KafkaMessage struct {
	ID      uint   `gorm:"primaryKey;column:id" json:"id"`
	Email   string `gorm:"primaryKey;column:email" json:"email"`
	Message string `gorm:"primaryKey;column:message" json:"message"`
}
