package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/anduckhmt146/kakfa-consumer/internal/db"
	router "github.com/anduckhmt146/kakfa-consumer/internal/routers"
	"github.com/anduckhmt146/kakfa-consumer/internal/services"
	"github.com/spf13/viper"

	"gorm.io/gorm"
)

type Server struct {
	DB         *gorm.DB
	httpServer *http.Server
}

func NewServer() *Server {
	initDB, err := db.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	return &Server{
		DB: initDB,
	}
}

func (s *Server) SetupHTTPServer() {
	httpPort := viper.GetString("service.http_port")
	ginRouter := router.SetupRouter(s.DB)

	s.httpServer = &http.Server{
		Addr:    httpPort,
		Handler: ginRouter,
	}

	go func() {
		log.Println("HTTP server listening on", httpPort)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve Gin HTTP server: %v", err)
		}
	}()
}

func (s *Server) Shutdown() {
	// Close database connection
	if sqlDB, err := s.DB.DB(); err != nil {
		log.Printf("Error on closing db connection: %v", err)
	} else {
		sqlDB.Close()
	}
	log.Println("Server and database connection closed")

}

func (s *Server) StartKafka() {
	// Kafka configuration
	brokerHost := viper.GetString("kafka.host")
	brokerPort := viper.GetString("kafka.port")
	groupID := viper.GetString("kafka.group_id")
	version := viper.GetString("kafka.version")
	topic := viper.GetString("kafka.topic")

	// Format brokers
	brokers := fmt.Sprintf("%s:%s", brokerHost, brokerPort)

	// Create a new Kafka consumer group
	consumerGroup := services.NewKafkaConsumerGroup(brokers, groupID, version)
	consumer := services.NewMessageConsumer()

	ctx := context.Background()
	topics := []string{topic}

	go func() {
		defer close(consumer.Ready)
		if err := consumerGroup.Client.Consume(ctx, topics, consumer); err != nil {
			log.Printf("Error from consumer: %v", err)
		}
	}()

	<-consumer.Ready
	log.Println("Kafka consumer started successfully!")

}

func (s *Server) Start() {
	// Start HTTP Server
	s.SetupHTTPServer()

	// Start Kafka
	go s.StartKafka()

	// Handle graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("Shutting down server...")
	s.Shutdown()
}
