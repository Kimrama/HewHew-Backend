package server

import (
	_shopController "hewhew-backend/pkg/shop/controller"
	_shopRepository "hewhew-backend/pkg/shop/repository"
	_shopService "hewhew-backend/pkg/shop/service"
    "hewhew-backend/utils"
)

func (s *fiberServer) initShopRouter() {
    shopRepository := _shopRepository.NewShopRepositoryImpl(s.db, s.conf.Supabase)
    shopService := _shopService.NewShopServiceImpl(shopRepository)
    shopController := _shopController.NewShopControllerImpl(shopService)


    canteenGroup := s.app.Group("/canteens")
    canteenGroup.Post("/", shopController.CreateCanteen)
    canteenGroup.Put("/:canteenName", shopController.EditCanteen)
    // canteenGroup.Delete("/:canteenName", shopController.DeleteCanteen)
    shopGroup := s.app.Group("/v1/shop")
    shopGroup.Use(utils.JWTProtected())
    shopGroup.Put("/", shopController.EditShop)
	shopGroup.Get("/", shopController.GetShop)
	shopGroup.Put("/toggle_open_state", shopController.ChangeState)
	shopGroup.Put("/shopimage", shopController.EditShopImage)
}