package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/menu/model"

	"github.com/google/uuid"
)

type MenuService interface {
	CreateMenu(menuModel *model.MenuRequest,shopID uuid.UUID) error
	GetMenusByShopID(shopID uuid.UUID) ([]*entities.Menu, error)
	DeleteMenu(menuID uuid.UUID, admin *entities.ShopAdmin) error
	EditMenu(menuID uuid.UUID, admin *entities.ShopAdmin,menuModel *model.MenuRequest) error
}
