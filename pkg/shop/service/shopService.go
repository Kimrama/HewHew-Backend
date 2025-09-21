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
	GetTagsByShopIDAndTopic(shopID string, topic string) ([]entities.Tag, error)
	EditTag(tagID string, topic string) error
	GetAllTags(shopID string) ([]entities.Tag, error)
	DeleteTag(tagID string) error
}
