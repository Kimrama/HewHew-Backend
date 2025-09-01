package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/model"
)

type UserService interface {
	CreateUser(userModel *model.CreateUserRequest) error
	GetUserByUsername(username string) (*entities.User, error)
	EditUser(userID string, user *entities.User) error
}
