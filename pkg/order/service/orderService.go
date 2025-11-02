package service

import (
	"hewhew-backend/pkg/order/model"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(orderModel *model.CreateOrderRequest, userID uuid.UUID) error
	AcceptOrder(acceptOrderModel *model.AcceptOrderRequest) error
	ConfirmOrder(confirmOrderModel *model.ConfirmOrderRequest, userID uuid.UUID) error
	DeleteOrder(orderID uuid.UUID, userID uuid.UUID) error
	GetOrdersByUserID(userID uuid.UUID) ([]model.GetOrderResponse, error)
	GetOrderByDeliveryUserID(userID uuid.UUID) ([]model.GetOrderResponse, error)
	GetOrdersByShopID(userID string) ([]model.GetOrderResponse, error)
	GetAvailableOrders() ([]model.GetAvailableOrderResponse, error)
	GetOrderByID(orderID uuid.UUID) (*model.GetOrderByIdResponse, error)
	GetUserAverageRating(userID uuid.UUID) (float64, error)

	CreateReview(reviewModel *model.CreateReviewRequest, userID uuid.UUID) error
	GetReviewsByTargetUserID(targetUserID uuid.UUID) ([]*model.GetReviewResponse, error)
	GetReviewsByReviewerUserID(reviewerUserID uuid.UUID) ([]*model.GetReviewResponse, error)
	GetReviewByID(reviewID uuid.UUID) (*model.GetReviewResponse, error)
	CreateTransactionLog(log *model.TransactionLog) error
	CreateNotification(notification *model.CreateNotificationRequest) error
	CreateNotificationDriver(notification *model.CreateNotificationDriverRequest) error
	GetNotificationByUserID(userID uuid.UUID) ([]*model.GetNotification, error)
}
