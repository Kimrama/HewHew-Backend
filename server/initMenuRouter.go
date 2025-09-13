
package server

import (
	_menuController "hewhew-backend/pkg/menu/controller"
	_menuRepository "hewhew-backend/pkg/menu/repository"
	_menuService "hewhew-backend/pkg/menu/service"
)

func (s *fiberServer) initMenuRouter() {
	menuRepository := _menuRepository.NewMenuRepositoryImpl(s.db)
	menuService := _menuService.NewMenuServiceImpl(menuRepository)
	menuController := _menuController.NewMenuControllerImpl(menuService)

	menuGroup := s.app.Group("/v1/menu")
	menuGroup.Post("/", menuController.CreateMenu)
}
