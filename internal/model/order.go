// internal/model/order.go
package model

import (
	"time"
)

type Order struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	ProductID  string    `json:"productId" gorm:"column:product_id;type:varchar(255);not null"`
	TotalPrice float64   `json:"totalPrice" gorm:"column:total_price;type:decimal(10,2);not null"`
	Status     string    `json:"status" gorm:"column:status;type:varchar(50);not null;default:'pending'"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (Order) TableName() string {
	return "orders"
}

type CreateOrderRequest struct {
	ProductID string `json:"productId" binding:"required" validate:"required,uuid4"`
	Quantity  int    `json:"quantity" binding:"required,min=1" validate:"required,min=1"`
}
