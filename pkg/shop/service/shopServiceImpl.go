package service

import (
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/shop/repository"
)

type ShopServiceImpl struct {
    ShopRepository repository.ShopRepository
}

func NewShopServiceImpl(ShopRepository repository.ShopRepository) ShopService {
    return &ShopServiceImpl{
        ShopRepository: ShopRepository,
    }
}
func (s *ShopServiceImpl) CreateCanteen(canteenModel interface{}) error {
    return s.ShopRepository.CreateCanteen(canteenModel)
}   

func (s *ShopServiceImpl) EditCanteen(canteenName string,canteenEntity *entities.Canteen) error {
    if canteenName == "" {
        return fmt.Errorf("canteen name is required")
    }
    return s.ShopRepository.EditCanteen(canteenName,canteenEntity)
}

func (s *ShopServiceImpl) DeleteCanteen(canteenID string) error {
    return nil
}
