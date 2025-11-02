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
	menuRepository := _menuRepository.NewMenuRepositoryImpl(s.db, s.conf.Supabase)
	menuService := _menuService.NewMenuServiceImpl(menuRepository)

	shopRepository := _shopRepository.NewShopRepositoryImpl(s.db, s.conf.Supabase)
	shopService := _shopService.NewShopServiceImpl(shopRepository)

	menuController := _menuController.NewMenuControllerImpl(menuService, shopService)

	menuGroup := s.app.Group("/v1/menu")

	menuGroup.Get("/:menu_id", menuController.GetMenuByID)
	menuGroup.Get("/menus/popular", menuController.PopularMenus)
	menuGroup.Use(utils.JWTProtected())
	menuGroup.Post("/", menuController.CreateMenu)
	menuGroup.Put("/:menu_id", menuController.EditMenu)
	menuGroup.Delete("/:menu_id", menuController.DeleteMenu)
	menuGroup.Get("/", menuController.GetAllMenu)
	menuGroup.Patch("/:menu_id/status", menuController.EditMenuStatus)
	menuGroup.Put("/:menu_id/image", menuController.EditMenuImage)

}
