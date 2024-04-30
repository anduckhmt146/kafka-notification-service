package router

import (
	"fmt"

	"github.com/anduckhmt146/kakfa-consumer/internal/handlers"
	"github.com/anduckhmt146/kakfa-consumer/internal/repositories"
	"github.com/anduckhmt146/kakfa-consumer/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	brokerHost := viper.GetString("kafka.host")
	brokerPort := viper.GetString("kafka.port")
	// Format brokers
	brokers := fmt.Sprintf("%s:%s", brokerHost, brokerPort)

	// User Repo
	userRepo := repositories.NewUserRepository(db)
	kafkaProducer, err := services.NewKafkaProducer(brokers)
	if err != nil {
		panic("Failed to create Kafka producer: " + err.Error())
	}

	// Notification Service
	notificationService := services.NewNotificationService(kafkaProducer.Producer, userRepo)

	notificationHandler := handlers.NewNotificationHandler(notificationService)
	heathHandler := handlers.NewHealthHandler()

	router.GET("/health", heathHandler.HealthCheck)
	router.POST("/message", notificationHandler.SendNotification)

	return router
}
