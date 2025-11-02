package controller

import "github.com/gofiber/fiber/v2"

type MenuController interface {
	CreateMenu(ctx *fiber.Ctx) error
	GetAllMenu(ctx *fiber.Ctx) error
	GetMenuByID(ctx *fiber.Ctx) error
	EditMenu(ctx *fiber.Ctx) error
	EditMenuStatus(ctx *fiber.Ctx) error
	EditMenuImage(ctx *fiber.Ctx) error
	DeleteMenu(ctx *fiber.Ctx) error
	PopularMenus(ctx *fiber.Ctx) error
}
