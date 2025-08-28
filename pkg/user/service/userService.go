package service

import "hewhew-backend/entities"

type UserService interface {
	CreateUser(userEntity *entities.User) (*entities.User, error)
	GetUsers() ([]*entities.User, error)
	UploadUserProfileImage(username string, ext string, image []byte) (string, error)
}
