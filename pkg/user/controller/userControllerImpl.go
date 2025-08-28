package controller

import (
	"fmt"
	"hewhew-backend/pkg/user/model"
	"hewhew-backend/pkg/user/service"
	"hewhew-backend/utils"

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
	UserName := ctx.FormValue("Username")
	Password := ctx.FormValue("Password")
	FName := ctx.FormValue("FName")
	LName := ctx.FormValue("LName")
	Gender := ctx.FormValue("Gender")
	image, err := ctx.FormFile("Image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to retrieve image",
		})
	}

	preprocessUploadImage, ext, err := utils.PreprocessUploadImage(image)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to preprocess image",
		})
	}
	url, err := c.userService.UploadUserProfileImage(UserName, ext, preprocessUploadImage)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload image",
		})
	}

	fmt.Println("Creating user:", UserName, Password, FName, LName, Gender, url)
	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *UserControllerImpl) GetUsers(ctx *fiber.Ctx) error {
	users, err := c.userService.GetUsers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}
	return ctx.JSON(users)
}

func (c *UserControllerImpl) DeleteUser(ctx *fiber.Ctx) error {
	return nil
}

func (c *UserControllerImpl) TestUser(ctx *fiber.Ctx) error {
	var user model.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}
	var res = model.Response{
		Name: user.Name,
	}
	return ctx.JSON(res)
}
