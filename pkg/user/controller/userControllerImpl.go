package controller

import (
	"hewhew-backend/pkg/user/service"

	"github.com/gofiber/fiber/v2"
)

type UserControllerImpl struct {
	userService service.UserService
}

func NewUserControllerImpl(userService service.UserService) UserController {
	return &UserControllerImpl{
		userService: userService,
	}
}

func (c *UserControllerImpl) CreateUser(ctx *fiber.Ctx) error {
	return nil
}
