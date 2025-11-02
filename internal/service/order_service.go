// internal/service/order_service.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"order-service/internal/model"
	"order-service/internal/rabbitmq"
	"order-service/internal/redis"
	"order-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(req *model.CreateOrderRequest) (*model.Order, error)
	GetOrdersByProductID(productID string) ([]model.Order, error)
}

type orderService struct {
	repo          repository.OrderRepository
	productClient ProductServiceClient
	redisClient   redis.ClientInterface
	publisher     rabbitmq.PublisherInterface
}

func NewOrderService(repo repository.OrderRepository, productClient ProductServiceClient, redisClient redis.ClientInterface, publisher rabbitmq.PublisherInterface) OrderService {
	return &orderService{
		repo:          repo,
		productClient: productClient,
		redisClient:   redisClient,
		publisher:     publisher,
	}
}

func (s *orderService) CreateOrder(req *model.CreateOrderRequest) (*model.Order, error) {
	product, err := s.productClient.GetProductByID(req.ProductID)
	log.Println(err)
	if err != nil {
		return nil, errors.New("failed to fetch product: " + err.Error())
	}
	totalPrice := product.Price * float64(req.Quantity)

	newOrder := &model.Order{
		ID:         uuid.New().String(),
		ProductID:  req.ProductID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	err = s.repo.Create(newOrder)
	if err != nil {
		return nil, err
	}

	// --- PUBLISH EVENT ---
	orderJSON, _ := json.Marshal(newOrder)
	err = s.publisher.Publish("amq.topic", "order.created", orderJSON)
	if err != nil {
		log.Printf("Failed to publish order.created event: %v", err)
	} else {
		log.Printf("Event 'order.created' published for order ID: %s", newOrder.ID)
	}

	return newOrder, nil
}

func (s *orderService) GetOrdersByProductID(productID string) ([]model.Order, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("orders:product:%s", productID)

	cachedData, err := s.redisClient.Get(ctx, cacheKey)
	if err == nil {
		log.Printf("Cache HIT for product %s", productID)
		var orders []model.Order
		if err := json.Unmarshal([]byte(cachedData), &orders); err == nil {
			return orders, nil
		}
	}
	log.Printf("Cache MISS for product %s. Fetching from DB.", productID)
	orders, err := s.repo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}

	if orders != nil {
		data, _ := json.Marshal(orders)
		s.redisClient.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	return orders, nil
}
