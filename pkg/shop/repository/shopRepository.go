package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type ShopRepository interface {
	CreateCanteen(canteenModel interface{}) error
	EditCanteen(canteenName string, canteenEntity *entities.Canteen) error
	DeleteCanteen(canteenName string) error
	EditShop(body entities.Shop, shop uuid.UUID) error
	ChangeState(state bool, shopID uuid.UUID) error
	EditShopImage(AdminID uuid.UUID, imageModel *utils.ImageModel) error
	UploadShopImage(shopID uuid.UUID, imageModel *utils.ImageModel) (string, error)
	GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error)
	GetShopAdminByUsername(username string) (*entities.ShopAdmin, error)
	GetAllCanteens() ([]entities.Canteen, error)
	GetCanteenByName(canteenName string) (*entities.Canteen, error)
	GetAllShops() ([]entities.Shop, error)
	GetShopByID(shopID uuid.UUID) (*entities.Shop, error)
	GetShopsByCanteens(canteenNames []string) ([]entities.Shop, error)
	CreateTag(tagModel *entities.Tag) (*entities.Tag, error)
	GetTagsByShopIDAndTopic(shopID string, topic string) ([]entities.Tag, error)
	EditTag(tagModel *entities.Tag) error
	GetAllTags(shopID uuid.UUID) ([]entities.Tag, error)
	DeleteTag(tagID string) error
	GetAllMenus(shopID uuid.UUID) ([]*entities.Menu, error)
	CreateTransactionLog(log *entities.TransactionLog) error
	CreateNotification(notification *entities.Notification) error
	GetDropOffByID(id uuid.UUID) (*entities.DropOffLocation, error)
	GetOrderByID(orderID uuid.UUID) (*entities.Order, error)
	GetMenuByID(menuID uuid.UUID) (*entities.Menu, error)
	GetTagByID(tagID uuid.UUID) (*entities.Tag, error)
	GetOrderIDsFromTransactionLog() ([]uuid.UUID, error)
	CountMenusFromOrders(orderIDs []uuid.UUID) (map[uuid.UUID]int, error)
	GetPopularShopsByMenuCounts(menuCounts map[uuid.UUID]int) ([]*entities.Shop, error)
}