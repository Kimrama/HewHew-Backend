package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/model"
	"hewhew-backend/utils"
)

type UserService interface {
	CreateUser(userModel *model.CreateUserRequest) error
	GetUserByUsername(username string) (*entities.User, error)
	EditUser(userID string, userEntity *entities.User) error
	EditUserProfileImage(userID string, imageModel *utils.ImageModel) error
}
