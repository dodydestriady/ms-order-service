// internal/service/order_service.go
package service

import (
	"errors"
	"log"
	"order-service/internal/model"
	"order-service/internal/repository"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(req *model.CreateOrderRequest) (*model.Order, error)
	GetOrdersByProductID(productID string) ([]model.Order, error)
}

type orderService struct {
	repo          repository.OrderRepository
	productClient ProductServiceClient
}

func NewOrderService(repo repository.OrderRepository, productClient ProductServiceClient) OrderService {
	return &orderService{
		repo:          repo,
		productClient: productClient,
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

	return newOrder, nil
}

func (s *orderService) GetOrdersByProductID(productID string) ([]model.Order, error) {
	return s.repo.GetByProductID(productID)
}
