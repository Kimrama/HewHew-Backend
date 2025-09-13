package controller

import (
	"hewhew-backend/pkg/menu/model"
	"hewhew-backend/pkg/menu/service"
	"hewhew-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type MenuControllerImpl struct {
    MenuService service.MenuService
}

func NewMenuControllerImpl(MenuService service.MenuService) MenuController {
    return &MenuControllerImpl{
        MenuService: MenuService,
    }
}

func (c *MenuControllerImpl) CreateMenu(ctx *fiber.Ctx) error {
	Name := ctx.FormValue("name")
	Detail := ctx.FormValue("detail")
	Price := ctx.FormValue("price")
	Status := ctx.FormValue("status") 
	Image, _ := ctx.FormFile("image")
	Tag1ID := ctx.FormValue("tag1_id")
	Tag2ID := ctx.FormValue("tag2_id")

	var imageModel *utils.ImageModel
	if Image != nil {
		preprocessUploadImage, err := utils.PreprocessUploadImage(Image)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to preprocess image",
			})
		}
		imageModel = preprocessUploadImage
	}

	menuStatus := model.MenuStatus(Status)
	if menuStatus != model.MenuStatusAvailable && menuStatus != model.MenuStatusUnavailable {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status, must be AVAILABLE or UNAVAILABLE",
		})
	}

	menuModel := &model.CreateMenuRequest{
		Name:   Name,
		Detail: Detail,
		Price:  Price,
		Status: menuStatus,
		Tag1ID: Tag1ID,
		Tag2ID: Tag2ID,
		Image:  imageModel,
	}

	err := c.MenuService.CreateMenu(menuModel)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Menu created successfully",
	})
}


func (c *MenuControllerImpl) GetMenu(ctx *fiber.Ctx) error {
    return nil
}

func (c *MenuControllerImpl) EditMenu(ctx *fiber.Ctx) error {
    return nil
}
