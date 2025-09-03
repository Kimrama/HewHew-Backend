package server

import (
	_userController "hewhew-backend/pkg/user/controller"
	_userRepository "hewhew-backend/pkg/user/repository"
	_userService "hewhew-backend/pkg/user/service"
)

func (s *fiberServer) initUserRouter() {
	userRepository := _userRepository.NewUserRepositoryImpl(s.db, s.conf.Supabase)
	userService := _userService.NewUserServiceImpl(userRepository)
	userController := _userController.NewUserControllerImpl(userService)

	userGroup := s.app.Group("/v1/user")
	userGroup.Get("/:username", userController.GetUserByUsername)
	userGroup.Post("/register", userController.CreateUser)
	userGroup.Post("/login", userController.LoginUser)
	userGroup.Put("/edit-profile-image/:username", userController.EditUserProfileImage)
	userGroup.Put("/edit-user/:userID", userController.EditUser)
}
