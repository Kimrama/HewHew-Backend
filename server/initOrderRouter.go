package server

import (
	_orderController "hewhew-backend/pkg/order/controller"
	_orderRepository "hewhew-backend/pkg/order/repository"
	_orderService "hewhew-backend/pkg/order/service"
	"hewhew-backend/utils"
)

func (s *fiberServer) initOrderRouter() {
	orderRepository := _orderRepository.NewOrderRepositoryImpl(s.db, s.conf.Supabase)
	orderService := _orderService.NewOrderServiceImpl(orderRepository)
	orderController := _orderController.NewOrderControllerImpl(orderService)

	orderGroup := s.app.Group("/v1/order")
	orderGroup.Get("/available", orderController.GetAvailableOrders)
	orderGroup.Get("/nearby", orderController.GetNearbyOrders)
	orderGroup.Post("/transaction_log", orderController.CreateTransactionLog)
	orderGroup.Use(utils.JWTProtected())
	orderGroup.Get("/notifications/:receiver_id", orderController.GetNotificationByUserID)
	orderGroup.Post("/", orderController.CreateOrder)
	orderGroup.Get("/", orderController.GetOrdersByUserID)
	orderGroup.Get("/delivery", orderController.GetOrderByDeliveryUserID)
	orderGroup.Get("/shop", orderController.GetOrdersByShopID)
	orderGroup.Get("/:id", orderController.GetOrderByID)
	orderGroup.Post("/accept", orderController.AcceptOrder)
	orderGroup.Post("/confirm", orderController.ConfirmOrder)
	// orderGroup.Delete("/delete/:id", orderController.DeleteOrder)

	reviewGroup := s.app.Group("/v1/review")
	reviewGroup.Use(utils.JWTProtected())
	reviewGroup.Get("/averagerating", orderController.GetUserAverageRating)
	reviewGroup.Get("/user/:targetUserID", orderController.GetReviewsByTargetUserID)
	reviewGroup.Get("/reviewer/:reviewerUserID", orderController.GetReviewsByReviewerUserID)
	reviewGroup.Get("/:reviewID", orderController.GetReviewByID)
	reviewGroup.Post("/", orderController.CreateReview)

}
