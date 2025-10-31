package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/order/model"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(orderModel *model.CreateOrderRequest, userID uuid.UUID) error
	AcceptOrder(acceptOrderModel *model.AcceptOrderRequest) error
	ConfirmOrder(confirmOrderModel *model.ConfirmOrderRequest, userID uuid.UUID) error
	DeleteOrder(orderID uuid.UUID, userID uuid.UUID) error
	GetOrdersByUserID(userID uuid.UUID) ([]*entities.Order, error)
	GetOrderByDeliveryUserID(userID uuid.UUID) ([]*entities.Order, error)
	GetOrdersByShopID(userID string) ([]*entities.Order, error)
	GetAvailableOrders() ([]*entities.Order, error)
	GetOrderByID(orderID uuid.UUID) (*model.OrderResponse, error)
	GetUserAverageRating(userID uuid.UUID) (float64, error)

	CreateReview(reviewModel *model.CreateReviewRequest, userID uuid.UUID) error
	GetReviewsByTargetUserID(targetUserID uuid.UUID) ([]*entities.Review, error)
	GetReviewByID(reviewID uuid.UUID) (*entities.Review, error)
}
