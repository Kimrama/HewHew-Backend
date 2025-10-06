package repository

import (
	"hewhew-backend/entities"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(orderEntity *entities.Order) error
	CreateMenuQuantity(menuQuantityEntity *entities.MenuQuantity) error
	// GetOrdersByUserID(userID string) (interface{}, error)
	// GetOrdersByShopID(shopID string) (interface{}, error)
	// UpdateOrderStatus(orderID string, status string) error
	GetMenuByID(menuID uuid.UUID) (*entities.Menu, error)
	GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOff, error)
}
