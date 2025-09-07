package server

import (
	_shopController "hewhew-backend/pkg/shop/controller"
	_shopRepository "hewhew-backend/pkg/shop/repository"
	_shopService "hewhew-backend/pkg/shop/service"
)

func (s *fiberServer) initShopRouter() {
    shopRepository := _shopRepository.NewShopRepositoryImpl(s.db, s.conf.Supabase)
    shopService := _shopService.NewShopServiceImpl(shopRepository)
    shopController := _shopController.NewShopControllerImpl(shopService)

    shopGroup := s.app.Group("/v1/shop")

    canteenGroup := shopGroup.Group("/canteens")
    canteenGroup.Post("/", shopController.CreateCanteen)
    canteenGroup.Put("/:canteenName", shopController.EditCanteen)
    // canteenGroup.Delete("/:canteenName", shopController.DeleteCanteen)
}