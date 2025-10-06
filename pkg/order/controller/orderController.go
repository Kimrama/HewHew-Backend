package controller

import "github.com/gofiber/fiber/v2"

type OrderController interface {
	CreateOrder(ctx *fiber.Ctx) error
	GetOrdersByUserID(ctx *fiber.Ctx) error
	GetOrdersByShopID(ctx *fiber.Ctx) error
	UpdateOrderStatus(ctx *fiber.Ctx) error
}
