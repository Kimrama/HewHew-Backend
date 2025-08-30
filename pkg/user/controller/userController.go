package controller

import "github.com/gofiber/fiber/v2"

type UserController interface {
	CreateUser(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error
	GetUserByUsername(ctx *fiber.Ctx) error
}
