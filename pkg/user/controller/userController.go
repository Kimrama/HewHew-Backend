package controller

import "github.com/gofiber/fiber/v2"

type UserController interface {
	CreateUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error
	GetUser(ctx *fiber.Ctx) error
	EditUser(ctx *fiber.Ctx) error
	EditUserProfileImage(ctx *fiber.Ctx) error
	LoginShopAdmin(ctx *fiber.Ctx) error
	CreateAdmin(ctx *fiber.Ctx) error
	EditShop(ctx *fiber.Ctx) error
	GetShop(ctx *fiber.Ctx) error
	ChangeState(ctx *fiber.Ctx) error
	EditShopImage(ctx *fiber.Ctx) error
}
