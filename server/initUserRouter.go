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
	userGroup.Get("/shop", userController.GetAllShops)
	userGroup.Post("/register", userController.CreateUser)
	userGroup.Post("/login", userController.LoginUser)
	userGroup.Get("/:userID", userController.GetUserByID)
	userGroup.Use(utils.JWTProtected())
	userGroup.Get("/", userController.GetUser)
	userGroup.Put("/profile_image", userController.EditUserProfileImage)
	userGroup.Put("/", userController.EditUser)
	userGroup.Put("/topup", userController.Topup)
	userGroup.Post("/contact", userController.CreateUserContact)
	userGroup.Delete("/contact/:contactID", userController.DeleteUserContact)

	adminGroup := s.app.Group("/v1/admin")
	adminGroup.Post("/login", userController.LoginShopAdmin)
	adminGroup.Post("/register", userController.CreateAdmin)
	adminGroup.Use(utils.JWTProtected())

}
