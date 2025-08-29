package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/model"
)

type UserService interface {
	CreateUser(userModel *model.UserModel) error
	GetUsers() ([]*entities.User, error)
}
