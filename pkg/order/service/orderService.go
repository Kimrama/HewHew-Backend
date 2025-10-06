package service

import (
	"hewhew-backend/pkg/order/model"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(orderModel *model.CreateOrderRequest, userID uuid.UUID) error
	// GetOrdersByUserID(userID string) (interface{}, error)
	// GetOrdersByShopID(shopID string) (interface{}, error)
	// UpdateOrderStatus(orderID string, status string) error
}
