package server

import (
	_shopController "hewhew-backend/pkg/shop/controller"
	_shopRepository "hewhew-backend/pkg/shop/repository"
	_shopService "hewhew-backend/pkg/shop/service"
)

func (s *fiberServer) initShopRouter() {
	shopRepository := _shopRepository.NewShopRepositoryImpl(s.db)
	shopService := _shopService.NewShopServiceImpl(shopRepository)
	shopController := _shopController.NewShopControllerImpl(shopService)

	shopGroup := s.app.Group("/v1/shop")
	shopGroup.Post("/canteens", shopController.CreateCanteen)
}