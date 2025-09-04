package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/model"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(userModel *model.CreateUserRequest) error
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByUserID(userID uuid.UUID) (*entities.User, error)
	EditUser(userID uuid.UUID, userEntity *entities.User) error
	EditUserProfileImage(userID uuid.UUID, imageModel *utils.ImageModel) error
}
