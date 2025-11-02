package controller

import (
	"hewhew-backend/pkg/menu/model"
	menusvc "hewhew-backend/pkg/menu/service"
	shopsvc "hewhew-backend/pkg/shop/service"
	"hewhew-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	menuModel := &model.MenuRequest{
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

func (c *MenuControllerImpl) GetAllMenu(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	isAdmin, ok := claims["admin"].(bool)
	if !ok || !isAdmin {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}
	shopIDStr, ok := claims["shop"].(string)
	if !ok || shopIDStr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid shop ID in token"})
	}

	menus, err := c.MenuService.GetMenusByShopID(shopID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(menus)
}

func (c *MenuControllerImpl) GetMenuByID(ctx *fiber.Ctx) error {
	menuIDStr := ctx.Params("menu_id")
	menuID, err := uuid.Parse(menuIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid menu ID",
		})
	}

	menuEntity, tags, err := c.MenuService.GetMenuByID(menuID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "menu not found",
		})
	}

	menu := &model.GetMenuByIDResponse{
		MenuID:   menuEntity.MenuID,
		Name:     menuEntity.Name,
		Detail:   menuEntity.Detail,
		Price:    menuEntity.Price,
		Status:   string(menuEntity.Status),
		ImageURL: menuEntity.ImageURL,
		Tags:     tags,
	}

	return ctx.JSON(menu)
}

func (c *MenuControllerImpl) DeleteMenu(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	menuIDParam := ctx.Params("menu_id")
	if menuIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "menuID is required"})
	}

	menuID, err := uuid.Parse(menuIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid menuID"})
	}

	admin, err := c.ShopService.GetShopAdminByUsername(username)
	if err != nil || admin == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Admin not found"})
	}

	if err := c.MenuService.DeleteMenu(menuID, admin); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Menu deleted successfully"})
}

func (c *MenuControllerImpl) EditMenu(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	menuIDParam := ctx.Params("menu_id")
	if menuIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "menuID is required"})
	}

	menuID, err := uuid.Parse(menuIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid menuID"})
	}

	admin, err := c.ShopService.GetShopAdminByUsername(username)
	if err != nil || admin == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Admin not found"})
	}

	var req model.EditMenuRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := c.MenuService.EditMenu(menuID, admin, &req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Menu updated successfully"})
}

func (c *MenuControllerImpl) EditMenuStatus(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	menuIDParam := ctx.Params("menu_id")
	if menuIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "menuID is required"})
	}
	menuID, err := uuid.Parse(menuIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid menuID"})
	}
	admin, err := c.ShopService.GetShopAdminByUsername(username)
	if err != nil || admin == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Admin not found"})
	}
	status := ctx.FormValue("status")
	if status != string(model.MenuStatusAvailable) && status != string(model.MenuStatusUnavailable) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid status"})
	}

	if err := c.MenuService.EditMenuStatus(menuID, admin, status); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"message": "Menu status updated successfully"})
}

func (c *MenuControllerImpl) EditMenuImage(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	admin, err := c.ShopService.GetShopAdminByUsername(username)
	if err != nil || admin == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Admin not found"})
	}

	menuIDParam := ctx.Params("menu_id")
	if menuIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "menuID is required"})
	}

	menuID, err := uuid.Parse(menuIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid menuID"})
	}
	Image, _ := ctx.FormFile("image")
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

	if err := c.MenuService.EditMenuImage(menuID, admin, imageModel); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"message": "Menu image updated successfully"})

}

func (c *MenuControllerImpl) PopularMenus(ctx *fiber.Ctx) error {
	menus, err := c.MenuService.GetPopularMenus()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(menus)
}
