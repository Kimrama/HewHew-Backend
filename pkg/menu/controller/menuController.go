package controller

import "github.com/gofiber/fiber/v2"

type MenuController interface {
	CreateMenu(ctx *fiber.Ctx) error
	GetMenu(ctx *fiber.Ctx) error
	EditMenu(ctx *fiber.Ctx) error
}
