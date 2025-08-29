package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/utils"
)

type UserRepository interface {
	CreateUser(userEntity *entities.User) error
	UploadUserProfileImage(username string, imageModel *utils.ImageModel) (string, error)
	GetUsers() ([]*entities.User, error)
}
