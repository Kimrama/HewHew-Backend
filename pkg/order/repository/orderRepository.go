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
	GetOrdersByUserID(userID uuid.UUID) ([]*entities.Order, error)
	GetOrderByDeliveryUserID(userID uuid.UUID) ([]*entities.Order, error)
	GetOrdersByShopID(shopID uuid.UUID) ([]*entities.Order, error)
	GetAvailableOrders() ([]*entities.Order, error)
	GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error)

	GetOrderByID(orderID uuid.UUID) (*entities.Order, error)
	GetMenuByID(menuID uuid.UUID) (*entities.Menu, error)
	GetShopByID(shopID uuid.UUID) (*entities.Shop, error)
	GetCanteenByName(name string) (*entities.Canteen, error)
	GetDropOffByID(id uuid.UUID) (*entities.DropOffLocation, error)

	CreateReview(reviewEntity *entities.Review) error
	GetReviewsByTargetUserID(targetUserID uuid.UUID) ([]*entities.Review, error)
	GetReviewsByReviewerUserID(reviewerUserID uuid.UUID) ([]*entities.Review, error)
	GetReviewByID(reviewID uuid.UUID) (*entities.Review, error)
	CreateTransactionLog(log *entities.TransactionLog) error
	CreateNotification(notification *entities.Notification) error
}
