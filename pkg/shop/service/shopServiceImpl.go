package service

import (
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

func (s *ShopServiceImpl) EditCanteen(canteenModel interface{}) error {
    return nil
}

func (s *ShopServiceImpl) DeleteCanteen(canteenID string) error {
    return nil
}

func (s *ShopServiceImpl) GetCanteens() (interface{}, error) {
    return nil, nil
}
