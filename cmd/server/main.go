package main

import (
	"log"
	"net/http"
	"order-service/config"
	"order-service/internal/consumers"
	"order-service/internal/database"
	"order-service/internal/handler"
	"order-service/internal/rabbitmq"
	"order-service/internal/redis"
	"order-service/internal/repository"
	"order-service/internal/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.InitConfig()

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	redisAddr := cfg.RedisHost + ":" + cfg.RedisPort
	redisClient := redis.NewClient(redisAddr)

	rabbitmqURL := cfg.RabbitMQURL
	publisher, err := rabbitmq.NewPublisher(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}
	defer publisher.Close()

	orderRepo := repository.NewOrderRepository(db)

	productServiceURL := cfg.ProductServiceURL
	productClient := service.NewProductServiceClient(productServiceURL, cfg)
	orderService := service.NewOrderService(orderRepo, productClient, redisClient, publisher)

	orderHandler := handler.NewOrderHandler(orderService)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	router.POST("/orders", orderHandler.CreateOrder)
	router.GET("/orders/product/:productId", orderHandler.GetOrdersByProductID)

	// RUNS rabbitmq
	go consumers.StartOrderLoggerConsumer(rabbitmqURL)
	go consumers.StartProductEventListener(rabbitmqURL)

	port := cfg.AppPort
	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to run server : %v", err)
		}
	}()

	log.Printf("Order service is running on port %s", port)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
