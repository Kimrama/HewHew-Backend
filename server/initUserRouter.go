package server

import (
	_userController "hewhew-backend/pkg/user/controller"
	_userRepository "hewhew-backend/pkg/user/repository"
	_userService "hewhew-backend/pkg/user/service"
)

func (s *fiberServer) InitUserRouter() {
	userRepository := _userRepository.NewUserRepositoryImpl(s.db)
	userService := _userService.NewUserServiceImpl(userRepository)
	userController := _userController.NewUserControllerImpl(userService)

	userGroup := s.app.Group("/users")
	userGroup.Post("/", userController.CreateUser)
}
