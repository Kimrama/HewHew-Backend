package controller

import (
	"github.com/gofiber/fiber/v2"
)

type OrderController interface {
	CreateOrder(ctx *fiber.Ctx) error
	AcceptOrder(ctx *fiber.Ctx) error
	ConfirmOrder(ctx *fiber.Ctx) error
	DeleteOrder(ctx *fiber.Ctx) error
	GetOrdersByUserID(ctx *fiber.Ctx) error  // for customer
	GetOrdersByShopID(ctx *fiber.Ctx) error  // for shop admin
	GetAvailableOrders(ctx *fiber.Ctx) error // for delivery user
	GetOrderByID(ctx *fiber.Ctx) error
	GetUserAverageRating(ctx *fiber.Ctx) error

	CreateReview(ctx *fiber.Ctx) error
	GetReviewsByTargetUserID(ctx *fiber.Ctx) error
	GetReviewByID(ctx *fiber.Ctx) error
}
