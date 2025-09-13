package repository

import (
	"github.com/google/uuid"
	"hewhew-backend/entities"
	"hewhew-backend/utils"
)

type ShopRepository interface {
	CreateCanteen(canteenModel interface{}) error
	EditCanteen(canteenName string, canteenEntity *entities.Canteen) error
	DeleteCanteen(canteenID string) error
	EditShop(body entities.Shop, shop uuid.UUID) error
	ChangeState(state bool, shopID uuid.UUID) error
	EditShopImage(AdminID uuid.UUID, imageModel *utils.ImageModel) error
	UploadShopImage(shopID uuid.UUID, imageModel *utils.ImageModel) (string, error)
	GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error)
	GetShopAdminByUsername(username string) (*entities.ShopAdmin, error)
  GetAllCanteens() ([]entities.Canteen, error)
  CreateTag(tagModel *entities.Tag) error
}
