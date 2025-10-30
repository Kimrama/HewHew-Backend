package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/order/model"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(orderEntity *entities.Order) error
	CreateMenuQuantity(menuQuantityEntity *entities.MenuQuantity) error
	AcceptOrder(acceptOrderModel *model.AcceptOrderRequest) error
	GetUserAverageRating(userID uuid.UUID) (float64, error)
	CountActiveOrdersByUser(userID uuid.UUID) (int64, error)
	ConfirmOrder(orderID uuid.UUID, imageUrl string) error
	DeleteOrder(orderID uuid.UUID) error
	UploadConfirmImage(orderID uuid.UUID, imageModel *utils.ImageModel) (string, error)
	GetOrderByID(orderID uuid.UUID) (*entities.Order, error)
	GetOrdersByUserID(userID uuid.UUID) ([]*entities.Order, error)
	GetOrdersByShopID(shopID uuid.UUID) ([]*entities.Order, error)
	GetAvailableOrders() ([]*entities.Order, error)
	GetMenuByID(menuID uuid.UUID) (*entities.Menu, error)
	GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error)
	GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOffLocation, error)
}
