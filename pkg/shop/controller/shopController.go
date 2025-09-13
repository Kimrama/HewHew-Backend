package controller

import "github.com/gofiber/fiber/v2"

type ShopController interface {
	CreateCanteen(ctx *fiber.Ctx) error
	EditCanteen(ctx *fiber.Ctx) error
	DeleteCanteen(ctx *fiber.Ctx) error
	EditShop(ctx *fiber.Ctx) error
	GetShop(ctx *fiber.Ctx) error
	ChangeState(ctx *fiber.Ctx) error
	EditShopImage(ctx *fiber.Ctx) error
	Createtag(ctx *fiber.Ctx) error

}
