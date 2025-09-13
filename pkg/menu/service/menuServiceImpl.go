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

func (s *MenuServiceImpl) CreateMenu(menuModel interface{}) error {
    // Implement the logic to create a menu item
    // This is a placeholder implementation 
    err := s.MenuRepository.CreateMenu(menuModel)
    if err != nil {
        return err
    }
    return nil
}