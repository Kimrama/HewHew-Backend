package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type MenuRepository interface {
	CreateMenu(menuEntity *entities.Menu) error
	UploadMenuImage(menuID uuid.UUID, imageModel *utils.ImageModel) (string, error)
	GetMenusByShopID(shopID uuid.UUID) ([]*entities.Menu, error)
	GetMenuByID(menuID uuid.UUID) (*entities.Menu, error)
	DeleteMenu(menuID uuid.UUID) error
	EditMenu(menuEntity *entities.Menu) error
	EditMenuStatus(menuID uuid.UUID, status string) error
	EditMenuImage(menuID uuid.UUID, imageModel *utils.ImageModel) error
	
}
