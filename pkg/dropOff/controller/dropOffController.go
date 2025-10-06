package controller

import "github.com/gofiber/fiber/v2"

type DropOffController interface {
	CreateDropOff(ctx *fiber.Ctx) error
	GetAllDropOffs(ctx *fiber.Ctx) error
	GetDropOffByID(ctx *fiber.Ctx) error
}
