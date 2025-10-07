package server

import (
	_orderController "hewhew-backend/pkg/order/controller"
	_orderRepository "hewhew-backend/pkg/order/repository"
	_orderService "hewhew-backend/pkg/order/service"
	"hewhew-backend/utils"
)

func (s *fiberServer) initOrderRouter() {
	orderRepository := _orderRepository.NewOrderRepositoryImpl(s.db)
	orderService := _orderService.NewOrderServiceImpl(orderRepository)
	orderController := _orderController.NewOrderControllerImpl(orderService)

	orderGroup := s.app.Group("/v1/order")
	orderGroup.Use(utils.JWTProtected())
	orderGroup.Post("/", orderController.CreateOrder)
	orderGroup.Get("/", orderController.GetOrdersByUserID)
	orderGroup.Get("/shop", orderController.GetOrdersByShopID)
	orderGroup.Get("/available", orderController.GetAvailableOrders)
	orderGroup.Get("/:id", orderController.GetOrderByID)
	orderGroup.Post("/accept", orderController.AcceptOrder)
	orderGroup.Post("/confirm/:id", orderController.ConfirmOrder)

}
