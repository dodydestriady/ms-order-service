package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"order-service/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *model.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByProductID(productID string) ([]model.Order, error) {
	args := m.Called(productID)

	ordersArg := args.Get(0)

	if ordersArg == nil {
		return nil, args.Error(1)
	}

	return ordersArg.([]model.Order), args.Error(1)
}

type MockProductServiceClient struct {
	mock.Mock
}

func (m *MockProductServiceClient) GetProductByID(productID string) (*model.Product, error) {
	args := m.Called(productID)

	productArg := args.Get(0)

	if productArg == nil {
		return nil, args.Error(1)
	}

	return productArg.(*model.Product), args.Error(1)
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

type MockRabbitMQPublisher struct {
	mock.Mock
}

func (m *MockRabbitMQPublisher) Publish(exchange, routingKey string, body []byte) error {
	args := m.Called(exchange, routingKey, body)
	return args.Error(0)
}

func (m *MockRabbitMQPublisher) Close() {
	m.Called()
}

// --- Unit Test ---

func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockProductClient := new(MockProductServiceClient)
	mockRedisClient := new(MockRedisClient)
	mockPublisher := new(MockRabbitMQPublisher)

	req := &model.CreateOrderRequest{ProductID: "prod-123", Quantity: 2}
	expectedProduct := &model.Product{ID: "prod-123", Price: 1000}
	expectedOrder := &model.Order{
		ID:         "order-456", // ID akan di-generate, kita asumsikan nilainya ini
		ProductID:  req.ProductID,
		TotalPrice: 2000,
		Status:     "pending",
	}

	mockProductClient.On("GetProductByID", req.ProductID).Return(expectedProduct, nil)
	mockRepo.On("Create", mock.AnythingOfType("*model.Order")).Return(nil).Run(func(args mock.Arguments) {
		order := args.Get(0).(*model.Order)
		order.ID = "order-456"
	})
	mockPublisher.On("Publish", "amq.topic", "order.created", mock.AnythingOfType("[]uint8")).Return(nil)

	service := &orderService{
		repo:          mockRepo,
		productClient: mockProductClient,
		redisClient:   mockRedisClient,
		publisher:     mockPublisher,
	}

	result, err := service.CreateOrder(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.ID, result.ID)
	assert.Equal(t, expectedOrder.ProductID, result.ProductID)
	assert.Equal(t, expectedOrder.TotalPrice, result.TotalPrice)
	assert.Equal(t, expectedOrder.Status, result.Status)

	mockProductClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestOrderService_CreateOrder_ProductNotFound(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockProductClient := new(MockProductServiceClient)
	mockRedisClient := new(MockRedisClient)
	mockPublisher := new(MockRabbitMQPublisher)

	req := &model.CreateOrderRequest{ProductID: "prod-999"}

	mockProductClient.On("GetProductByID", req.ProductID).Return(nil, errors.New("product not found"))

	service := &orderService{
		repo:          mockRepo,
		productClient: mockProductClient,
		redisClient:   mockRedisClient,
		publisher:     mockPublisher,
	}

	result, err := service.CreateOrder(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to fetch product")

	mockRepo.AssertNotCalled(t, "Create")
	mockPublisher.AssertNotCalled(t, "Publish")
	mockProductClient.AssertExpectations(t)
}
