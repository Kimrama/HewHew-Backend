package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/menu/model"
	"hewhew-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MenuService interface {
	CreateMenu(menuModel *model.MenuRequest, shopID uuid.UUID) error
	GetMenusByShopID(shopID uuid.UUID) ([]*entities.Menu, error)
	GetMenuByID(menuID uuid.UUID) (*entities.Menu, []string, error)
	DeleteMenu(menuID uuid.UUID, admin *entities.ShopAdmin) error
	EditMenu(menuID uuid.UUID, admin *entities.ShopAdmin, menuModel *model.MenuRequest) error
	EditMenuStatus(menuID uuid.UUID, admin *entities.ShopAdmin, status string) error
	EditMenuImage(menuID uuid.UUID, admin *entities.ShopAdmin, imageModel *utils.ImageModel) error
	GetPopularMenus() (fiber.Map, error)
}
