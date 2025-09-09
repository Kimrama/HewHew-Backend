package controller

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/shop/model"
	"hewhew-backend/pkg/shop/service"
	"net/url"

	"github.com/gofiber/fiber/v2"
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
		Latitude:  body.Latitude,
		Longitude:  body.Longitude,
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
    canteenName := ctx.Params("canteenName")
    decodedName, err := url.PathUnescape(canteenName)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid canteen name",
        })
    }
    if err := s.ShopService.DeleteCanteen(decodedName); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return ctx.JSON(fiber.Map{"message": "Canteen deleted successfully"})
}

