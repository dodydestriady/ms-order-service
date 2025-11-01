package repository

import "order-service/internal/model"

type OrderRepository interface {
	Create(order *model.Order) error
	GetByProductID(productID string) ([]model.Order, error)
}
