package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(userEntity *entities.User) error
	UploadUserProfileImage(username string, imageModel *utils.ImageModel) (string, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByUserID(userID uuid.UUID) (*entities.User, error)
	EditUser(userID uuid.UUID, user *entities.User) error
	EditUserProfileImage(userID uuid.UUID, imageModel *utils.ImageModel) error
}
