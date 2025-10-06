package server

import (
	_orderController "hewhew-backend/pkg/order/controller"
	_orderRepository "hewhew-backend/pkg/order/repository"
	_orderService "hewhew-backend/pkg/order/service"
)

func (s *fiberServer) initOrderRouter() {
	orderRepository := _orderRepository.NewOrderRepositoryImpl(s.db)
	orderService := _orderService.NewOrderServiceImpl(orderRepository)
	orderController := _orderController.NewOrderControllerImpl(orderService)

	orderGroup := s.app.Group("/v1/order")
	orderGroup.Post("/", orderController.CreateOrder)
	orderGroup.Get("/", orderController.GetOrdersByUserID)
	orderGroup.Get("/shop", orderController.GetOrdersByShopID)
	orderGroup.Put("/:order_id/status", orderController.UpdateOrderStatus)
	
}
