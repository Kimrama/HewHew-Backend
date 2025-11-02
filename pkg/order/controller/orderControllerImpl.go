package controller

import (
	"hewhew-backend/pkg/order/model"
	"hewhew-backend/pkg/order/service"
	"hewhew-backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderControllerImpl struct {
	OrderService service.OrderService
}

func NewOrderControllerImpl(OrderService service.OrderService) OrderController {
	return &OrderControllerImpl{
		OrderService: OrderService,
	}
}

func (oc *OrderControllerImpl) CreateOrder(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	userIDstr, ok := claims["user_id"].(string)
	if !ok || userIDstr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID in token"})
	}

	var req model.CreateOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = oc.OrderService.CreateOrder(&req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
	})
}

func (oc *OrderControllerImpl) AcceptOrder(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	userIDstr, ok := claims["user_id"].(string)
	if !ok || userIDstr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID in token"})
	}

	orderIDstr := ctx.FormValue("order_id")
	orderID, err := uuid.Parse(orderIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order ID"})
	}

	acceptOrderModel := &model.AcceptOrderRequest{
		DeliveryuserID: userID,
		OrderID:        orderID,
	}

	err = oc.OrderService.AcceptOrder(acceptOrderModel)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order accepted successfully",
	})

}

func (oc *OrderControllerImpl) ConfirmOrder(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	userIDstr, ok := claims["user_id"].(string)
	if !ok || userIDstr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID in token"})
	}

	orderIDstr := ctx.FormValue("order_id")
	Image, _ := ctx.FormFile("image")

	orderID, err := uuid.Parse(orderIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order ID"})
	}

	var imageModel *utils.ImageModel
	if Image != nil {
		preprocessUploadImage, err := utils.PreprocessUploadImage(Image)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to preprocess image",
			})
		}
		imageModel = preprocessUploadImage
	}

	ConfirmOrderModel := &model.ConfirmOrderRequest{
		OrderID: orderID,
		Image:   imageModel,
	}

	err = oc.OrderService.ConfirmOrder(ConfirmOrderModel, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "order confirmed successfully",
	})
}

func (oc *OrderControllerImpl) DeleteOrder(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	userIDstr, ok := claims["user_id"].(string)
	if !ok || userIDstr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID in token"})
	}

	orderIDstr := ctx.Params("id")
	orderID, err := uuid.Parse(orderIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order ID"})
	}

	if err := oc.OrderService.DeleteOrder(orderID, userID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Menu deleted successfully"})
}

func (oc *OrderControllerImpl) GetOrdersByShopID(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tokenUserID, ok := claims["user_id"].(string)
	if !ok || tokenUserID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	orders, err := oc.OrderService.GetOrdersByShopID(tokenUserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return ctx.JSON(orders)
}

func (oc *OrderControllerImpl) GetOrdersByUserID(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tokenUserID, ok := claims["user_id"].(string)
	if !ok || tokenUserID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	userID, err := uuid.Parse(tokenUserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	orders, err := oc.OrderService.GetOrdersByUserID(userID)
	if err != nil {
		return err
	}
	return ctx.JSON(orders)
}

func (oc *OrderControllerImpl) GetAvailableOrders(ctx *fiber.Ctx) error {
	orders, err := oc.OrderService.GetAvailableOrders()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch available orders"})
	}
	return ctx.JSON(orders)
}

func (oc *OrderControllerImpl) GetOrderByDeliveryUserID(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	tokenUserID, ok := claims["user_id"].(string)
	if !ok || tokenUserID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	userID, err := uuid.Parse(tokenUserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}
	orders, err := oc.OrderService.GetOrderByDeliveryUserID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}
	return ctx.JSON(orders)
}

func (oc *OrderControllerImpl) GetOrderByID(ctx *fiber.Ctx) error {
	orderIDstr := ctx.Params("id")
	orderID, err := uuid.Parse(orderIDstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order ID"})
	}
	orderModel, err := oc.OrderService.GetOrderByID(orderID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if orderModel == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "order not found"})
	}
	return ctx.JSON(orderModel)
}

func (oc *OrderControllerImpl) GetNearbyOrders(ctx *fiber.Ctx) error {
	var body model.GetNearbyOrdersRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	lat, err := strconv.ParseFloat(body.Latitude, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid latitude",
		})
	}
	lon, err := strconv.ParseFloat(body.Longitude, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid longitude",
		})
	}

	shops, err := oc.OrderService.GetNearbyOrders(lat, lon)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(shops)
}

func (oc *OrderControllerImpl) GetUserAverageRating(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	tokenUserID, ok := claims["user_id"].(string)
	if !ok || tokenUserID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	userID, err := uuid.Parse(tokenUserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}
	averageRating, err := oc.OrderService.GetUserAverageRating(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch average rating"})
	}
	return ctx.JSON(fiber.Map{"average_rating": averageRating})
}

func (oc *OrderControllerImpl) CreateReview(ctx *fiber.Ctx) error {
	claims, err := utils.GetClaimsFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tokenUserID, ok := claims["user_id"].(string)
	if !ok || tokenUserID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	userID, err := uuid.Parse(tokenUserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	var reviewModel model.CreateReviewRequest
	if err := ctx.BodyParser(&reviewModel); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	err = oc.OrderService.CreateReview(&reviewModel, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Review created successfully",
	})
}

func (oc *OrderControllerImpl) GetReviewsByTargetUserID(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("targetUserID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	reviews, err := oc.OrderService.GetReviewsByTargetUserID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(reviews)
}

func (oc *OrderControllerImpl) GetReviewsByReviewerUserID(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("reviewerUserID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}
	reviews, err := oc.OrderService.GetReviewsByReviewerUserID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(reviews)
}

func (oc *OrderControllerImpl) GetReviewByID(ctx *fiber.Ctx) error {
	reviewIDStr := ctx.Params("reviewID")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid review ID"})
	}

	review, err := oc.OrderService.GetReviewByID(reviewID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if review == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "review not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(review)
}

func (oc *OrderControllerImpl) CreateTransactionLog(ctx *fiber.Ctx) error {
	var body model.TransactionLog
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	// ส่ง model ให้ service แปลงเป็น entity เอง
	if err := oc.OrderService.CreateTransactionLog(&body); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Transaction log created successfully",
	})
}

func (oc *OrderControllerImpl) CreateNotification(ctx *fiber.Ctx) error {
	var body model.CreateNotificationRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	if err := oc.OrderService.CreateNotification(&body); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Notification created successfully",
	})
}
