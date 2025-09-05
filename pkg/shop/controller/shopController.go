package controller

import "github.com/gofiber/fiber/v2"

type ShopController interface {
	CreateCanteen(ctx *fiber.Ctx) error
}
