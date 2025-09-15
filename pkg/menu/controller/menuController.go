package controller

import "github.com/gofiber/fiber/v2"

type MenuController interface {
	CreateMenu(ctx *fiber.Ctx) error
	GetAllMenu(ctx *fiber.Ctx) error
	EditMenu(ctx *fiber.Ctx) error
	DeleteMenu(ctx *fiber.Ctx) error
}
