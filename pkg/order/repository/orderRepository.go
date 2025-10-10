package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/order/model"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(orderEntity *entities.Order) error
	CreateMenuQuantity(menuQuantityEntity *entities.MenuQuantity) error
	AcceptOrder(acceptOrderModel *model.AcceptOrderRequest) error
	GetOrderByID(orderID uuid.UUID) (*entities.Order, error)
	GetOrdersByUserID(userID uuid.UUID) ([]*entities.Order, error)
	GetOrdersByShopID(shopID uuid.UUID) ([]*entities.Order, error)
	GetAvailableOrders() ([]*entities.Order, error)
	GetMenuByID(menuID uuid.UUID) (*entities.Menu, error)
	GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error)
	GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOffLocation, error)
}
