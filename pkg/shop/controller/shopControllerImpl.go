package controller

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/shop/model"
	"hewhew-backend/pkg/shop/service"

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
    name := ctx.FormValue("Name")
    latitude := ctx.FormValue("Latitude")
    longitude := ctx.FormValue("Longitude")

    canteenModel := &model.CreateCanteenRequest{
        CanteenName:   name,
        Latitude:  latitude,
        Longitude: longitude,
    }
    if err := s.ShopService.CreateCanteen(canteenModel); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Canteen created successfully",
    })
}

func (s *ShopControllerImpl) EditCanteen(ctx *fiber.Ctx) error {
    var body model.EditcanteenRequest
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

	if err := c.userService.EditCanteen(c); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

    return ctx.JSON(fiber.Map{"message": " Canteen updated successfully "})

}

func (s *ShopControllerImpl) DeleteCanteen(ctx *fiber.Ctx) error {
    return nil
}

func (s *ShopControllerImpl) GetCanteens(ctx *fiber.Ctx) error {
    return nil
}
