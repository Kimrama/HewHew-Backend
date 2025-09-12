package controller

import (
	"errors"
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/shop/model"
	"hewhew-backend/pkg/shop/service"
	"hewhew-backend/utils"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopControllerImpl struct {
	ShopService service.ShopService
}

func NewShopControllerImpl(ShopService service.ShopService) ShopController {
	return &ShopControllerImpl{
		ShopService: ShopService,
	}
}

func (s *ShopControllerImpl) CreateCanteen(ctx *fiber.Ctx) error {
	var body model.CanteenRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	c := &entities.Canteen{
		CanteenName: body.CanteenName,
		Latitude:    body.Latitude,
		Longitude:   body.Longitude,
	}
	if err := s.ShopService.CreateCanteen(c); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Canteen created successfully",
	})
}

func (s *ShopControllerImpl) EditCanteen(ctx *fiber.Ctx) error {
	canteenName := ctx.Params("canteenName")
	decodedName, err := url.PathUnescape(canteenName)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid canteen name",
		})
	}

	var body model.CanteenRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	c := &entities.Canteen{
		CanteenName: decodedName,
		Latitude:    body.Latitude,
		Longitude:   body.Longitude,
	}

	if err := s.ShopService.EditCanteen(decodedName, c); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{"message": "Canteen updated successfully"})
}

func (s *ShopControllerImpl) DeleteCanteen(ctx *fiber.Ctx) error {
	return nil
}

func (c *ShopControllerImpl) ChangeState(ctx *fiber.Ctx) error {
	var body model.ChangeState
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
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

	admin, err := c.ShopService.GetShopAdminByUsername(claims["username"].(string))
	if err != nil || admin == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Admin not found",
		})
	}

	if err := c.ShopService.ChangeState(body, admin.ShopID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{"message": "Shop state updated successfully"})
}

func (c *ShopControllerImpl) EditShop(ctx *fiber.Ctx) error {
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

	var body model.EditShopRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}
	admin, err := c.ShopService.GetShopAdminByUsername(claims["username"].(string))

	if err != nil || admin == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Admin not found",
		})
	}

	if err := c.ShopService.EditShop(body, admin.ShopID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{"message": "Shop updated successfully"})
}

func (c *ShopControllerImpl) GetShop(ctx *fiber.Ctx) error {
	claims, err := getClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	userIDStr, _ := claims["user_id"].(string)
	if userIDStr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user token required"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id in token"})
	}

	fmt.Println("userID from token:", userID)

	shop, err := c.ShopService.GetShopByAdminID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "shop not found"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"name":         shop.Name,
		"canteen_name": shop.CanteenName,
		"state":        shop.State,
	})
}

func (c *ShopControllerImpl) EditShopImage(ctx *fiber.Ctx) error {

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

	tokenUserID, ok := claims["user_id"].(string)
	if !ok || tokenUserID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	fmt.Println("userID from token:", tokenUserID)
	shop, err := c.ShopService.GetShopByAdminID(uuid.MustParse(tokenUserID))
	fmt.Println("Shop from DB:", shop)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("Shop ID from DB:", shop.ShopID)

	imageModel, err := utils.PreprocessUploadImage(image)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to preprocess image",
		})
	}

	adminUUID, err := uuid.Parse(tokenUserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid admin id in token",
		})
	}

	err = c.ShopService.EditShopImage(adminUUID, imageModel)
	fmt.Print("After EditShopImage")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Shop image updated successfully"})
}

func getClaimsFromToken(ctx *fiber.Ctx) (jwt.MapClaims, error) {
	token := ctx.Locals("jwt").(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}
	return claims, nil
}

func (c *ShopControllerImpl) GetAllCanteens(ctx *fiber.Ctx) error {
	canteens, err := c.ShopService.GetAllCanteens()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"canteens": canteens,
	})
}
