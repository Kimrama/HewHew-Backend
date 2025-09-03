package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(userEntity *entities.User) error
	UploadUserProfileImage(username string, imageModel *utils.ImageModel) (string, error)
 	EditUserProfileImage(username string, imageModel *utils.ImageModel) (string, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByUserID(userID uuid.UUID) (*entities.User, error)
	EditUser(userID string, user *entities.User) error
}
