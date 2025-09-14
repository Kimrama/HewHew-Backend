package server

import (
	_menuController "hewhew-backend/pkg/menu/controller"
	_menuRepository "hewhew-backend/pkg/menu/repository"
	_menuService "hewhew-backend/pkg/menu/service"

	_shopRepository "hewhew-backend/pkg/shop/repository"
	_shopService "hewhew-backend/pkg/shop/service"

	"hewhew-backend/utils"
)

func (s *fiberServer) initMenuRouter() {
	// menu
	menuRepository := _menuRepository.NewMenuRepositoryImpl(s.db, s.conf.Supabase)
	menuService := _menuService.NewMenuServiceImpl(menuRepository)

	// shop
	shopRepository := _shopRepository.NewShopRepositoryImpl(s.db, s.conf.Supabase)
	shopService := _shopService.NewShopServiceImpl(shopRepository)

	// controller needs both
	menuController := _menuController.NewMenuControllerImpl(menuService, shopService)

	menuGroup := s.app.Group("/v1/menu")
	menuGroup.Use(utils.JWTProtected())
	menuGroup.Post("/", menuController.CreateMenu)
	menuGroup.Put("/:menu_id", menuController.EditMenu)
	menuGroup.Delete("/:menu_id", menuController.DeleteMenu)
	menuGroup.Get("/", menuController.GetAllMenu)

}
