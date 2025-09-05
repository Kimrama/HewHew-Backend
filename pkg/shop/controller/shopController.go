package controller

import "github.com/gofiber/fiber/v2"

type ShopController interface {
	createShop(ctx *fiber.Ctx) error
}
