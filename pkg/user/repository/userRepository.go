package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/utils"
)

type UserRepository interface {
	CreateUser(userEntity *entities.User) error
	UploadUserProfileImage(username string, imageModel *utils.ImageModel) (string, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByUserID(userID string) (*entities.User, error)
	EditUser(userID string, user *entities.User) error
	EditUserProfileImage(userID string, imageModel *utils.ImageModel) error
}
