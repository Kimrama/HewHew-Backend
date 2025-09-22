package controller

import (
	"errors"
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/model"
	"hewhew-backend/pkg/user/service"
	"hewhew-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	password := ctx.FormValue("password")
	username := ctx.FormValue("username")
	fname := ctx.FormValue("fname")
	lname := ctx.FormValue("lname")
	gender := ctx.FormValue("gender")
	image, _ := ctx.FormFile("image")

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
	fmt.Println(userModel.Username)

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
			"error": "Invalid username or password" + err.Error(),
		})
	}
	if !utils.CompareHashPassword(user.Password, loginRequest.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	token, err := utils.GenerateUserJWT(user.UserID, user.Username)
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

func (c *UserControllerImpl) GetShop(ctx *fiber.Ctx) error {
	claims, err := getClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	userIDStr, _ := claims["user_id"].(string) // ✅ ใช้ user_id
	if userIDStr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user token required"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id in token"})
	}

	fmt.Println("userID from token:", userID)

	shop, err := c.userService.GetShopByAdminID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "shop not found"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("Shop Name:", shop)
	fmt.Println("Canteen Name:", shop.CanteenName)

	return ctx.JSON(fiber.Map{
		"name":         shop.Name,
		"canteen_name": shop.CanteenName,
		"shopimg":      shop.ImageURL,
		"state":       shop.State,
	})
}

func (c *UserControllerImpl) EditUser(ctx *fiber.Ctx) error {
	// ใช้ helper
	claims, err := getClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tokenUserID, ok := claims["user_id"].(string)
	if !ok || tokenUserID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	var body model.EditUserRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	userUUID, err := uuid.Parse(tokenUserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID in token",
		})
	}

	u := &entities.User{
		UserID: userUUID,
		FName:  body.FName,
		LName:  body.LName,
		Gender: body.Gender,
	}

	if err := c.userService.EditUser(userUUID, u); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{"message": "User updated successfully"})
}

func (c *UserControllerImpl) CreateAdmin(ctx *fiber.Ctx) error {
	password := ctx.FormValue("password")
	username := ctx.FormValue("username")
	fname := ctx.FormValue("fname")
	lname := ctx.FormValue("lname")

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	AdminModel := &model.CreateAdminRequest{
		Username: username,
		Password: hashedPassword,
		FName:    fname,
		LName:    lname,
	}
	fmt.Println(AdminModel.Username)

	if err := c.userService.CreateAdmin(AdminModel); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func (c *UserControllerImpl) LoginShopAdmin(ctx *fiber.Ctx) error {
	var req model.ShopAdminLoginRequest_and_Shop
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	admin, err := c.userService.GetShopAdminByUsername(req.Username)
	if err != nil || admin == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	if !utils.CompareHashPassword(admin.Password, req.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	token, err := utils.GenerateAdminJWT(admin.AdminID, admin.ShopID, admin.Username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (c *UserControllerImpl) Topup(ctx *fiber.Ctx) error {

	var body model.TopupRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	if body.Amount <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "amount must be greater than zero",
		})
	}

	claims, err := getClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tokenUserID, ok := claims["user_id"].(string)
	if !ok || tokenUserID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	err = c.userService.Topup(tokenUserID, body.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{"message": "Topup successfully"})
}
