package service

import (
	"hewhew-backend/pkg/menu/repository"
)

type MenuServiceImpl struct {
    MenuRepository repository.MenuRepository
}

func NewMenuServiceImpl(MenuRepository repository.MenuRepository) MenuService {
    return &MenuServiceImpl{
        MenuRepository: MenuRepository,
    }
}



