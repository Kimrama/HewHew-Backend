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
	GetTagByID(tagID uuid.UUID) (*entities.Tag, error)
	DeleteMenu(menuID uuid.UUID) error
	EditMenu(menuID uuid.UUID, updates map[string]interface{}) error
	EditMenuStatus(menuID uuid.UUID, status string) error
	EditMenuImage(menuID uuid.UUID, imageModel *utils.ImageModel) error
	GetOrderIDsFromTransactionLog() ([]uuid.UUID, error)
	CountMenusFromOrders(orderIDs []uuid.UUID) (map[uuid.UUID]int, error)
	GetMenusByIDs(menuCounts map[uuid.UUID]int) ([]*entities.Menu, error)
}
