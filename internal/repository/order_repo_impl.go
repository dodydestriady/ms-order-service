// internal/repository/order_repo_impl.go
package repository

import (
	"order-service/internal/model"

	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}

func (r *orderRepositoryImpl) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepositoryImpl) GetByProductID(productID string) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Where("product_id = ?", productID).Order("created_at desc").Find(&orders).Error
	return orders, err
}
