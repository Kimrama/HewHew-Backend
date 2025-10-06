package service

import (
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/order/model"
	"hewhew-backend/pkg/order/repository"
	"time"

	"github.com/google/uuid"
)

type OrderServiceImpl struct {
	OrderRepository repository.OrderRepository
}

func NewOrderServiceImpl(OrderRepository repository.OrderRepository) OrderService {
	return &OrderServiceImpl{
		OrderRepository: OrderRepository,
	}
}

func (os *OrderServiceImpl) CreateOrder(orderModel *model.CreateOrderRequest, userID uuid.UUID) error {
	for _, item := range orderModel.Menu {
		menu, err := os.OrderRepository.GetMenuByID(item.MenuID)
		if err != nil {
			return fmt.Errorf("menu with ID %s not found", item.MenuID)
		}
		if menu.Status != "available" {
			return fmt.Errorf("menu %s is not available", menu.Name)
		}
	}

	_, err := os.OrderRepository.GetDropOffByID(orderModel.DropOffID)
	if err != nil {
		return fmt.Errorf("dropoff location with ID %s not found", orderModel.DropOffID)
	}

	OrderEntity := &entities.Order{
		OrderID:         uuid.New(),
		UserOrderID:     userID,
		UserDeliveryID:  nil,
		Status:          "Un Paid",
		OrderDate:       time.Now(),
		DeliveryMethod:  string(orderModel.DeliveryMethod),
		AppointmentTime: orderModel.AppointmentTime,
		DropOffID:       orderModel.DropOffID,
	}

	if err := os.OrderRepository.CreateOrder(OrderEntity); err != nil {
		return err
	}

	for _, item := range orderModel.Menu {
		menuQuantityEntity := &entities.MenuQuantity{
			MenuQuantityID: uuid.New(),
			MenuID:         item.MenuID,
			OrderID:        OrderEntity.OrderID,
			Quantity:       item.Quantity,
		}

		if err := os.OrderRepository.CreateMenuQuantity(menuQuantityEntity); err != nil {
			return err
		}
	}

	return nil
}

// func (os *OrderServiceImpl) GetOrdersByUserID(userID string) (interface{}, error) {
// 	return os.OrderRepository.GetOrdersByUserID(userID)
// }

// func (os *OrderServiceImpl) GetOrdersByShopID(shopID string) (interface{}, error) {
// 	return os.OrderRepository.GetOrdersByShopID(shopID)
// }

// func (os *OrderServiceImpl) UpdateOrderStatus(orderID string, status string) error {
// 	return os.OrderRepository.UpdateOrderStatus(orderID, status)
// }
