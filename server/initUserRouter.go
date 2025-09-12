package server

import (
	_userController "hewhew-backend/pkg/user/controller"
	_userRepository "hewhew-backend/pkg/user/repository"
	_userService "hewhew-backend/pkg/user/service"
	"hewhew-backend/utils"
)

func (s *fiberServer) initUserRouter() {
	userRepository := _userRepository.NewUserRepositoryImpl(s.db, s.conf.Supabase)
	userService := _userService.NewUserServiceImpl(userRepository)
	userController := _userController.NewUserControllerImpl(userService)
	userGroup := s.app.Group("/v1/user")
	userGroup.Post("/register", userController.CreateUser)
	userGroup.Post("/login", userController.LoginUser)
	userGroup.Use(utils.JWTProtected())
	userGroup.Get("/", userController.GetUser)
	userGroup.Put("/profile-image", userController.EditUserProfileImage)
	userGroup.Put("/", userController.EditUser)


	adminGroup := s.app.Group("/v1/admin")
	adminGroup.Post("/login", userController.LoginShopAdmin)
	adminGroup.Post("/register", userController.CreateAdmin)
	adminGroup.Use(utils.JWTProtected())


	

}
