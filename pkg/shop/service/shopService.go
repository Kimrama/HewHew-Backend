package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/shop/model"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type ShopService interface {
	CreateCanteen(canteenModel interface{}) error
	EditCanteen(canteenName string, canteenEntity *entities.Canteen) error
	DeleteCanteen(canteenName string) error
	GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error)
	EditShop(body model.EditShopRequest, shop uuid.UUID) error
	ChangeState(body model.ChangeState, admin_id uuid.UUID) error
	EditShopImage(shopID uuid.UUID, imageModel *utils.ImageModel) error
	GetShopAdminByUsername(username string) (*entities.ShopAdmin, error)
	CreateTag(ShopID string, body *model.TagCreateRequest) (*entities.Tag, error)
	GetAllCanteens() ([]entities.Canteen, error)
	GetCanteenByName(canteenName string) (*entities.Canteen, error)
	GetAllShops() ([]entities.Shop, error)
	GetShopByID(shopID uuid.UUID) (*entities.Shop, error)
	GetNearbyShops(lat, lon float64) ([]model.GetNearByShopResponse, error)
	GetTagsByShopIDAndTopic(shopID string, topic string) ([]entities.Tag, error)
	EditTag(tagID string, topic string) error
	GetAllTags(shopID uuid.UUID) ([]entities.Tag, error)
	DeleteTag(tagID string) error
	GetAllMenus(shopID uuid.UUID) ([]*model.GetMenuByIDResponse, error)
	GetPopularShops() ([]*entities.Shop, error)
}
