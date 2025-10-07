package controller

import (
	"hewhew-backend/pkg/dropOff/model"
	"hewhew-backend/pkg/dropOff/service"
	"hewhew-backend/utils"

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
	latitude := ctx.FormValue("latitude")
	longitude := ctx.FormValue("longitude")
	name := ctx.FormValue("name")
	detail := ctx.FormValue("detail")
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

	dropoffModel := &model.CreateDropOffRequest{
		Latitude:  latitude,
		Longitude: longitude,
		Name:      name,
		Detail:    detail,
		Image:     imageModel,
	}

	if err := dc.DropOffService.CreateDropOff(dropoffModel); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Drop off location created successfully",
	})
}

func (dc *DropOffControllerImpl) GetAllDropOffs(ctx *fiber.Ctx) error {
	dropoffEntities, err := dc.DropOffService.GetAllDropOffs()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var dropoffs []model.DropOffDetailResponse
	for _, entity := range dropoffEntities {
		dropoff := model.DropOffDetailResponse{
			Latitude:  entity.Latitude,
			Longitude: entity.Longitude,
			Name:      entity.Name,
			Detail:    entity.Detail,
			Image:     entity.ImageURL,
		}
		dropoffs = append(dropoffs, dropoff)
	}

	return ctx.JSON(dropoffs)
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
