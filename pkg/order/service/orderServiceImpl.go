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
		OrderID:           uuid.New(),
		UserOrderID:       userID,
		UserDeliveryID:    nil,
		Status:            "waiting",
		OrderDate:         time.Now(),
		DeliveryMethod:    string(orderModel.DeliveryMethod),
		AppointmentTime:   orderModel.AppointmentTime,
		DropOffLocationID: orderModel.DropOffID,
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

func (os *OrderServiceImpl) AcceptOrder(acceptOrderModel *model.AcceptOrderRequest) error {
	order, err := os.OrderRepository.GetOrderByID(acceptOrderModel.OrderID)
	if err != nil {
		return fmt.Errorf("order with ID %s not found", acceptOrderModel.OrderID)
	}
	if order.Status != "waiting" {
		return fmt.Errorf("order with ID %s is not in a state to be accepted", acceptOrderModel.OrderID)
	}

	if err := os.OrderRepository.AcceptOrder(acceptOrderModel); err != nil {
		return err
	}
	return nil
}

func (os *OrderServiceImpl) GetOrdersByUserID(userID uuid.UUID) ([]*entities.Order, error) {
	return os.OrderRepository.GetOrdersByUserID(userID)
}

func (os *OrderServiceImpl) GetOrdersByShopID(userID string) ([]*entities.Order, error) {
	adminID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	admin, err := os.OrderRepository.GetShopByAdminID(adminID)
	if err != nil {
		return nil, fmt.Errorf("admin not found: %w", err)
	}
	orders, err := os.OrderRepository.GetOrdersByShopID(admin.ShopID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	return orders, nil
}

func (os *OrderServiceImpl) GetAvailableOrders() ([]*entities.Order, error) {
	orders, err := os.OrderRepository.GetAvailableOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	return orders, nil
}

func (os *OrderServiceImpl) GetOrderByID(orderID uuid.UUID) (*entities.Order, error) {
	order, err := os.OrderRepository.GetOrderByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	return order, nil
}
