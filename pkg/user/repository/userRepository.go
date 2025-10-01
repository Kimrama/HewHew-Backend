package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(userEntity *entities.User) error
	CreateAdmin(adminModel *entities.ShopAdmin) error
	CreateShop(shopEntity *entities.Shop) error
	UploadUserProfileImage(username string, imageModel *utils.ImageModel) (string, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByUserID(userID uuid.UUID) (*entities.User, error)
	EditUser(userID uuid.UUID, user *entities.User) error
	EditUserProfileImage(userID uuid.UUID, imageModel *utils.ImageModel) error
	GetShopAdminByUsername(username string) (*entities.ShopAdmin, error)
	GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error)
	Topup(topupModel *entities.TopUp) error
	GetAllShops() ([]entities.Shop, error)
}


