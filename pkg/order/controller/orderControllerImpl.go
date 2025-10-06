package controller

import (
	"hewhew-backend/pkg/order/model"
	"hewhew-backend/pkg/order/service"
	"hewhew-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	userIDstr, ok := claims["user_id"].(string)
	if !ok || userIDstr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID in token"})
	}

	var req model.CreateOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = oc.OrderService.CreateOrder(&req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
	})
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
