package controller

import (
	"hewhew-backend/pkg/dropOff/model"
	"hewhew-backend/pkg/dropOff/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DropOffControllerImpl struct {
	DropOffService service.DropOffService
}

func NewDropOffControllerImpl(DropOffService service.DropOffService) DropOffController {
	return &DropOffControllerImpl{
		DropOffService: DropOffService,
	}
}

func (dc *DropOffControllerImpl) CreateDropOff(ctx *fiber.Ctx) error {
	var model model.DropOffRequest
	if err := ctx.BodyParser(&model); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if err := dc.DropOffService.CreateDropOff(&model); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Drop off location created successfully",
	})
}

func (dc *DropOffControllerImpl) GetAllDropOffs(ctx *fiber.Ctx) error {
	dropOffs, err := dc.DropOffService.GetAllDropOffs()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(dropOffs)
}

func (dc *DropOffControllerImpl) GetDropOffByID(ctx *fiber.Ctx) error {
	dropOffIDstr := ctx.Params("dropOffID")
	if dropOffIDstr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "dropOffID is required",
		})
	}
	dropOffID, err := uuid.Parse(dropOffIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid dropOffID",
		})
	}

	dropOff, err := dc.DropOffService.GetDropOffByID(dropOffID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(dropOff)
}
