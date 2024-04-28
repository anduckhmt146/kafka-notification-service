package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/anduckhmt146/kakfa-consumer/internal/db"
	"github.com/anduckhmt146/kakfa-consumer/internal/services"
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/gorm"
)

type Server struct {
	DB         *gorm.DB
	httpServer *http.Server
	grpcServer *grpc.Server
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

func (s *Server) SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	return router
}

func (s *Server) SetupHTTPServer() {
	httpPort := viper.GetString("service.http_port")
	router := s.SetupRouter()

	s.httpServer = &http.Server{
		Addr:    httpPort,
		Handler: router,
	}

	go func() {
		log.Println("HTTP server listening on", httpPort)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve Gin HTTP server: %v", err)
		}
	}()
}

func (s *Server) SetupGRPCServer() {
	grpcPort := viper.GetString("service.grpc_port")
	grpcListener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer()),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer()),
	)

	healthService := &HealthService{}
	healthpb.RegisterHealthServer(s.grpcServer, healthService)
	grpc_prometheus.Register(s.grpcServer)
	grpc_prometheus.EnableHandlingTimeHistogram()

	go func() {
		log.Println("gRPC server listening on", grpcPort)
		if err := s.grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()
}

func (s *Server) Shutdown() {
	// Graceful shutdown for gRPC server
	s.grpcServer.GracefulStop()

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

	// Start gRPC Server
	s.SetupGRPCServer()

	// Handle graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("Shutting down server...")
	s.Shutdown()
}

type HealthService struct{}

func (s *HealthService) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthService) Watch(req *healthpb.HealthCheckRequest, server healthpb.Health_WatchServer) error {
	return nil
}
