package repository

import "hewhew-backend/entities"

type UserRepository interface {
	CreateUser(userEntity *entities.User) (*entities.User, error)
}
