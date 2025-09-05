
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
