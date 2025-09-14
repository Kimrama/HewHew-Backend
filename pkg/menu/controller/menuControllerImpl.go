package controller

import (
	"hewhew-backend/pkg/menu/model"
	menusvc "hewhew-backend/pkg/menu/service"
	shopsvc "hewhew-backend/pkg/shop/service"
	"hewhew-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type MenuControllerImpl struct {
    MenuService menusvc.MenuService
    ShopService shopsvc.ShopService
}

func NewMenuControllerImpl(MenuService menusvc.MenuService, shopService shopsvc.ShopService) MenuController {
    return &MenuControllerImpl{
        MenuService: MenuService,
        ShopService: shopService,
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

    claims, err := utils.GetClaimsFromToken(ctx)
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

    status := Status
    if status == "" {
        status = "unavailable" 
    }
    menuStatus := model.MenuStatus(status)
    if menuStatus != model.MenuStatusAvailable && menuStatus != model.MenuStatusUnavailable {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid status, must be available or unavailable",
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

	err = c.MenuService.CreateMenu(menuModel, admin.ShopID)
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

