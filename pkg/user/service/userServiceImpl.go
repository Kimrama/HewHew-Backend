package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/repository"
)

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserServiceImpl(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

func (s *UserServiceImpl) CreateUser(userEntity *entities.User) (*entities.User, error) {
	return s.userRepository.CreateUser(userEntity)
}
