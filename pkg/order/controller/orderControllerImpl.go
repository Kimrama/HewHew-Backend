package controller

import (
	"hewhew-backend/pkg/order/service"

	"github.com/gofiber/fiber/v2"
)

type OrderControllerImpl struct {
    OrderService service.OrderService
}

func NewOrderControllerImpl(OrderService service.OrderService) OrderController {
    return &OrderControllerImpl{
        OrderService: OrderService,
    }
}

func (oc *OrderControllerImpl) CreateOrder(ctx *fiber.Ctx) error {
    return nil
}

func (oc *OrderControllerImpl) GetOrdersByUserID(ctx *fiber.Ctx) error {
    return nil
}   

func (oc *OrderControllerImpl) GetOrdersByShopID(ctx *fiber.Ctx) error {
    return nil
}

func (oc *OrderControllerImpl) UpdateOrderStatus(ctx *fiber.Ctx) error {
    return nil
}