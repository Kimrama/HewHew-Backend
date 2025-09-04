package controller

import (
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/model"
	"hewhew-backend/pkg/user/service"
	"hewhew-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	password := ctx.FormValue("Password")
	username := ctx.FormValue("Username")
	fname := ctx.FormValue("FName")
	lname := ctx.FormValue("LName")
	gender := ctx.FormValue("Gender")
	image, _ := ctx.FormFile("Image")

	var imageModel *utils.ImageModel
	if image != nil {
		preprocessUploadImage, err := utils.PreprocessUploadImage(image)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to preprocess image",
			})
		}
		imageModel = preprocessUploadImage
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	userModel := &model.CreateUserRequest{
		Username: username,
		Password: hashedPassword,
		FName:    fname,
		LName:    lname,
		Gender:   gender,
		Image:    imageModel,
	}

	if err := c.userService.CreateUser(userModel); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func (c *UserControllerImpl) EditUserProfileImage(ctx *fiber.Ctx) error {
	image, err := ctx.FormFile("Image")
	if err != nil || image == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image file is required",
		})
	}

	claims, err := getClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userID := claims["user_id"].(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID in token",
		})
	}
	imageModel, err := utils.PreprocessUploadImage(image)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to preprocess image",
		})
	}

	err = c.userService.EditUserProfileImage(userUUID, imageModel)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profile image updated successfully"})
}

func (c *UserControllerImpl) LoginUser(ctx *fiber.Ctx) error {
	var loginRequest model.LoginRequest
	if err := ctx.BodyParser(&loginRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	user, err := c.userService.GetUserByUsername(loginRequest.Username)
	if err != nil || user == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	if !utils.CompareHashPassword(user.Password, loginRequest.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	token, err := utils.GenerateJWT(user.UserID, user.Username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}
	return ctx.JSON(fiber.Map{
		"token": token,
	})
}
func getClaimsFromToken(ctx *fiber.Ctx) (jwt.MapClaims, error) {
	token := ctx.Locals("jwt").(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}
	return claims, nil
}

func (c *UserControllerImpl) GetUser(ctx *fiber.Ctx) error {

	claims, err := getClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userID := claims["user_id"].(string)

	fmt.Println("userID from token:", userID)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID in token",
		})
	}
	userEntity, err := c.userService.GetUserByUserID(userUUID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}
	user := &model.UserDetailResponse{
		Username:        userEntity.Username,
		FName:           userEntity.FName,
		LName:           userEntity.LName,
		Gender:          userEntity.Gender,
		ProfileImageURL: userEntity.ProfileImageURL,
	}
	return ctx.JSON(user)
}

func (c *UserControllerImpl) EditUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	var body model.EditUserRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	body.UserID = userUUID

	u := &entities.User{
		UserID: userUUID,
		FName:  body.FName,
		LName:  body.LName,
		Gender: body.Gender,
	}

	if err := c.userService.EditUser(u.UserID, u); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "User updated successfully"})
}
