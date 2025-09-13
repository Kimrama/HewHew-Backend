package controller

import (
	"hewhew-backend/pkg/menu/service"
)

type MenuControllerImpl struct {
    MenuService service.MenuService
}

func NewMenuControllerImpl(MenuService service.MenuService) MenuController {
    return &MenuControllerImpl{
        MenuService: MenuService,
    }
}

