package service

import "hewhew-backend/entities"

type UserService interface {
	CreateUser(userEntity *entities.User) (*entities.User, error)
}
