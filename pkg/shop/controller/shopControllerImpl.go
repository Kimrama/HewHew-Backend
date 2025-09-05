package controller

import (
	"hewhew-backend/pkg/shop/service"

	"github.com/gofiber/fiber/v2"
)

type ShopControllerImpl struct {
    ShopService service.ShopService
}

func NewShopControllerImpl(ShopService service.ShopService) ShopController {
    return &ShopControllerImpl{
        ShopService: ShopService,
    }
}

func (s *ShopControllerImpl) CreateCanteen(ctx *fiber.Ctx) error {
    return nil
    
}
