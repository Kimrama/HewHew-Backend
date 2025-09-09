package controller

import "github.com/gofiber/fiber/v2"

type ShopController interface {
	CreateCanteen(ctx *fiber.Ctx) error
	EditCanteen(ctx *fiber.Ctx) error
	DeleteCanteen(ctx *fiber.Ctx) error
}
