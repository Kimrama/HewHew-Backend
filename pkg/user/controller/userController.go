package controller

import "github.com/gofiber/fiber/v2"

type UserController interface {
	CreateUser(ctx *fiber.Ctx) error
	GetUsers(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error
	TestUser(ctx *fiber.Ctx) error
}
