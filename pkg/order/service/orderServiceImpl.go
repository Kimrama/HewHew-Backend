package service

import (
	"errors"
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/order/model"
	"hewhew-backend/pkg/order/repository"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type OrderServiceImpl struct {
	OrderRepository repository.OrderRepository
}

func NewOrderServiceImpl(OrderRepository repository.OrderRepository) OrderService {
	return &OrderServiceImpl{
		OrderRepository: OrderRepository,
	}
}

func (os *OrderServiceImpl) CreateOrder(orderModel *model.CreateOrderRequest, userID uuid.UUID) error {
	var canteenName string

	for _, item := range orderModel.Menu {
		menu, err := os.OrderRepository.GetMenuByID(item.MenuID)
		if err != nil {
			return fmt.Errorf("menu with ID %s not found", item.MenuID)
		}
		if menu.Status != "available" {
			return fmt.Errorf("menu %s is not available", menu.Name)
		}

		shop, err := os.OrderRepository.GetShopByID(menu.ShopID)
		if err != nil {
			return fmt.Errorf("shop with ID %s not found for menu %s", menu.ShopID, menu.Name)
		}

		if canteenName == "" {
			canteenName = shop.CanteenName
		} else if shop.CanteenName != canteenName {
			return fmt.Errorf("all menus must be from the same canteen; found %s and %s", canteenName, shop.CanteenName)
		}
	}

	_, err := os.OrderRepository.GetDropOffByID(orderModel.DropOffID)
	if err != nil {
		return fmt.Errorf("dropoff location with ID %s not found", orderModel.DropOffID)
	}

	OrderEntity := &entities.Order{
		OrderID:           uuid.New(),
		UserOrderID:       userID,
		UserDeliveryID:    nil,
		Status:            "waiting",
		OrderDate:         time.Now(),
		DeliveryMethod:    string(orderModel.DeliveryMethod),
		AppointmentTime:   orderModel.AppointmentTime,
		DropOffLocationID: orderModel.DropOffID,
	}

	if err := os.OrderRepository.CreateOrder(OrderEntity); err != nil {
		return err
	}

	for _, item := range orderModel.Menu {
		if item.Quantity <= 0 {
			return fmt.Errorf("invalid quantity %d for menu ID %s", item.Quantity, item.MenuID)
		}
		menuQuantityEntity := &entities.MenuQuantity{
			MenuQuantityID: uuid.New(),
			MenuID:         item.MenuID,
			OrderID:        OrderEntity.OrderID,
			Quantity:       item.Quantity,
		}

		if err := os.OrderRepository.CreateMenuQuantity(menuQuantityEntity); err != nil {
			return err
		}
	}

	return nil
}

func (os *OrderServiceImpl) AcceptOrder(acceptOrderModel *model.AcceptOrderRequest) error {
	rating, err := os.OrderRepository.GetUserAverageRating(acceptOrderModel.DeliveryuserID)
	if err != nil {
		return fmt.Errorf("failed to fetch user rating: %v", err)
	}

	var maxOrders int
	switch {
	case rating == 0:
		maxOrders = 1
	case rating < 3.5:
		maxOrders = 1
	case rating < 4.0:
		maxOrders = 2
	case rating < 4.5:
		maxOrders = 3
	default:
		maxOrders = 4
	}

	count, err := os.OrderRepository.CountActiveOrdersByUser(acceptOrderModel.DeliveryuserID)
	if err != nil {
		return fmt.Errorf("failed to count active orders: %v", err)
	}
	if count >= int64(maxOrders) {
		return fmt.Errorf("you have reached the maximum allowed active orders (%d)", maxOrders)
	}

	order, err := os.OrderRepository.GetOrderByID(acceptOrderModel.OrderID)
	if err != nil {
		return fmt.Errorf("order with ID %s not found", acceptOrderModel.OrderID)
	}
	if order.Status != "waiting" {
		return fmt.Errorf("order with ID %s is not in a state to be accepted", acceptOrderModel.OrderID)
	}

	if err := os.OrderRepository.AcceptOrder(acceptOrderModel); err != nil {
		return err
	}
	return nil
}

func (os *OrderServiceImpl) ConfirmOrder(confirmOrderModel *model.ConfirmOrderRequest, userID uuid.UUID) error {
	order, err := os.OrderRepository.GetOrderByID(confirmOrderModel.OrderID)
	if err != nil {
		return fmt.Errorf("order with ID %s not found", confirmOrderModel.OrderID)
	}
	if *order.UserDeliveryID != userID {
		return fmt.Errorf("unauthorized: user does not own this order")
	}
	if order.Status != "accepted" {
		return fmt.Errorf("order with ID %s is not in a state to be confirm", confirmOrderModel.OrderID)
	}

	imageUrl := ""
	if confirmOrderModel.Image != nil {
		var err error
		imageUrl, err = os.OrderRepository.UploadConfirmImage(confirmOrderModel.OrderID, confirmOrderModel.Image)
		if err != nil {
			return err
		}
	}

	if err := os.OrderRepository.ConfirmOrder(confirmOrderModel.OrderID, imageUrl); err != nil {
		return err
	}
	return nil

}

func (os *OrderServiceImpl) DeleteOrder(orderID uuid.UUID, userID uuid.UUID) error {
	order, err := os.OrderRepository.GetOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	if order.UserOrderID != userID {
		return fmt.Errorf("unauthorized to delete this order")
	}

	return os.OrderRepository.DeleteOrder(orderID)
}

func (os *OrderServiceImpl) GetOrdersByUserID(userID uuid.UUID) ([]model.GetOrderResponse, error) {
	orders, err := os.OrderRepository.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []model.GetOrderResponse
	for _, order := range orders {
		resp, err := os.buildOrderResponse(order)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *resp)
	}
	return responses, nil
}

func (os *OrderServiceImpl) GetOrdersByShopID(userID string) ([]model.GetOrderResponse, error) {
	adminID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	admin, err := os.OrderRepository.GetShopByAdminID(adminID)
	if err != nil {
		return nil, fmt.Errorf("admin not found: %w", err)
	}

	orders, err := os.OrderRepository.GetOrdersByShopID(admin.ShopID)
	if err != nil {
		return nil, err
	}

	var responses []model.GetOrderResponse
	for _, order := range orders {
		resp, err := os.buildOrderResponse(order)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *resp)
	}
	return responses, nil
}

func (os *OrderServiceImpl) GetOrderByDeliveryUserID(userID uuid.UUID) ([]model.GetOrderResponse, error) {
	orders, err := os.OrderRepository.GetOrderByDeliveryUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []model.GetOrderResponse
	for _, order := range orders {
		resp, err := os.buildOrderResponse(order)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *resp)
	}
	return responses, nil
}

func (os *OrderServiceImpl) GetAvailableOrders() ([]model.GetAvailableOrderResponse, error) {
	orders, err := os.OrderRepository.GetAvailableOrders()
	if err != nil {
		return nil, err
	}

	responses := make([]model.GetAvailableOrderResponse, 0, len(orders))

	for _, order := range orders {
		if len(order.MenuQuantity) == 0 {
			continue
		}

		firstMenuID := order.MenuQuantity[0].MenuID
		menu, err := os.OrderRepository.GetMenuByID(firstMenuID)
		if err != nil {
			continue
		}

		shop, err := os.OrderRepository.GetShopByID(menu.ShopID)
		if err != nil {
			continue
		}

		canteen, err := os.OrderRepository.GetCanteenByName(shop.CanteenName)
		if err != nil {
			continue
		}

		dropOff, err := os.OrderRepository.GetDropOffByID(order.DropOffLocationID)
		if err != nil {
			continue
		}

		cLat, err1 := strconv.ParseFloat(canteen.Latitude, 64)
		cLon, err2 := strconv.ParseFloat(canteen.Longitude, 64)
		dLat, err3 := strconv.ParseFloat(dropOff.Latitude, 64)
		dLon, err4 := strconv.ParseFloat(dropOff.Longitude, 64)
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			continue
		}

		distance := calculateDistance(cLat, cLon, dLat, dLon)
		shippingFee := calculateShippingFee(distance)
		var amount float64
		menuQuantityResp := make([]model.MenuQuantityResponse, 0, len(order.MenuQuantity))
		for _, mq := range order.MenuQuantity {
			menu, err := os.OrderRepository.GetMenuByID(mq.MenuID)
			if err != nil {
				return nil, err
			}
			amount += menu.Price * float64(mq.Quantity)
			menuQuantityResp = append(menuQuantityResp, model.MenuQuantityResponse{
				MenuID:   mq.MenuID,
				Quantity: mq.Quantity,
			})
		}

		resp := model.GetAvailableOrderResponse{
			OrderID:           order.OrderID,
			UserOrderID:       order.UserOrderID,
			Status:            order.Status,
			OrderDate:         order.OrderDate,
			DeliveryMethod:    order.DeliveryMethod,
			AppointmentTime:   order.AppointmentTime,
			DropOffLocationID: order.DropOffLocationID,
			MenuQuantity:      menuQuantityResp,
			ShopName:          shop.Name,
			CanteenName:       canteen.CanteenName,
			ShippingFee:       shippingFee,
			Amount:            amount,
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (os *OrderServiceImpl) GetOrderByID(orderID uuid.UUID) (*model.GetOrderByIdResponse, error) {
	order, err := os.OrderRepository.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil || len(order.MenuQuantity) == 0 {
		return nil, errors.New("order not found or has no menu items")
	}

	firstMenu, err := os.OrderRepository.GetMenuByID(order.MenuQuantity[0].MenuID)
	if err != nil {
		return nil, err
	}

	shop, err := os.OrderRepository.GetShopByID(firstMenu.ShopID)
	if err != nil {
		return nil, err
	}

	canteen, err := os.OrderRepository.GetCanteenByName(shop.CanteenName)
	if err != nil {
		return nil, err
	}

	dropOff, err := os.OrderRepository.GetDropOffByID(order.DropOffLocationID)
	if err != nil {
		return nil, err
	}

	cLat, _ := strconv.ParseFloat(canteen.Latitude, 64)
	cLon, _ := strconv.ParseFloat(canteen.Longitude, 64)
	dLat, _ := strconv.ParseFloat(dropOff.Latitude, 64)
	dLon, _ := strconv.ParseFloat(dropOff.Longitude, 64)

	distance := calculateDistance(cLat, cLon, dLat, dLon)
	shippingFee := calculateShippingFee(distance)

	var amount float64
	var menuQuantityResp []model.MenuQuantityResponse
	for _, mq := range order.MenuQuantity {
		menu, err := os.OrderRepository.GetMenuByID(mq.MenuID)
		if err != nil {
			return nil, err
		}
		amount += menu.Price * float64(mq.Quantity)
		menuQuantityResp = append(menuQuantityResp, model.MenuQuantityResponse{
			MenuID:   mq.MenuID,
			Quantity: mq.Quantity,
		})
	}

	return &model.GetOrderByIdResponse{
		OrderID:              order.OrderID,
		UserOrderID:          order.UserOrderID,
		UserDeliveryID:       order.UserDeliveryID,
		Status:               order.Status,
		OrderDate:            order.OrderDate,
		DeliveryMethod:       order.DeliveryMethod,
		ConfirmationImageURL: order.ConfirmationImageURL,
		AppointmentTime:      order.AppointmentTime,
		DropOffLocationID:    order.DropOffLocationID,
		MenuQuantity:         menuQuantityResp,
		ShopName:             shop.Name,
		CanteenName:          canteen.CanteenName,
		ShippingFee:          shippingFee,
		Amount:               amount,
	}, nil
}

func (os *OrderServiceImpl) GetUserAverageRating(userID uuid.UUID) (float64, error) {
	return os.OrderRepository.GetUserAverageRating(userID)
}

func (os *OrderServiceImpl) CreateReview(reviewModel *model.CreateReviewRequest, userID uuid.UUID) error {
	reviewEntity := &entities.Review{
		ReviewID:       uuid.New(),
		UserReviewerID: userID,
		UserTargetID:   reviewModel.UserTargetID,
		OrderID:        reviewModel.OrderID,
		Rating:         reviewModel.Rating,
		Comment:        reviewModel.Comment,
		TimeStamp:      time.Now(),
	}

	if err := os.OrderRepository.CreateReview(reviewEntity); err != nil {
		return err
	}
	return nil
}

func (os *OrderServiceImpl) GetReviewsByTargetUserID(userID uuid.UUID) ([]*model.GetReviewResponse, error) {
	reviews, err := os.OrderRepository.GetReviewsByTargetUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []*model.GetReviewResponse
	for _, r := range reviews {
		response = append(response, &model.GetReviewResponse{
			ReviewID:       r.ReviewID,
			UserReviewerID: r.UserReviewerID,
			UserTargetID:   r.UserTargetID,
			OrderID:        r.OrderID,
			Rating:         r.Rating,
			Comment:        r.Comment,
			TimeStamp:      r.TimeStamp,
		})
	}

	return response, nil
}

func (os *OrderServiceImpl) GetReviewsByReviewerUserID(userID uuid.UUID) ([]*model.GetReviewResponse, error) {
	reviews, err := os.OrderRepository.GetReviewsByReviewerUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []*model.GetReviewResponse
	for _, r := range reviews {
		response = append(response, &model.GetReviewResponse{
			ReviewID:       r.ReviewID,
			UserReviewerID: r.UserReviewerID,
			UserTargetID:   r.UserTargetID,
			OrderID:        r.OrderID,
			Rating:         r.Rating,
			Comment:        r.Comment,
			TimeStamp:      r.TimeStamp,
		})
	}
	return response, nil
}

func (os *OrderServiceImpl) GetReviewByID(reviewID uuid.UUID) (*model.GetReviewResponse, error) {
	review, err := os.OrderRepository.GetReviewByID(reviewID)
	if err != nil || review == nil {
		return nil, err
	}

	return &model.GetReviewResponse{
		ReviewID:       review.ReviewID,
		UserReviewerID: review.UserReviewerID,
		UserTargetID:   review.UserTargetID,
		OrderID:        review.OrderID,
		Rating:         review.Rating,
		Comment:        review.Comment,
		TimeStamp:      review.TimeStamp,
	}, nil
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func calculateShippingFee(distanceKm float64) float64 {
	return math.Round(distanceKm*10*100) / 100
}

func (os *OrderServiceImpl) buildOrderResponse(order *entities.Order) (*model.GetOrderResponse, error) {
	if order == nil || len(order.MenuQuantity) == 0 {
		return nil, errors.New("order not found or has no menu items")
	}

	firstMenu, err := os.OrderRepository.GetMenuByID(order.MenuQuantity[0].MenuID)
	if err != nil {
		return nil, err
	}

	shop, err := os.OrderRepository.GetShopByID(firstMenu.ShopID)
	if err != nil {
		return nil, err
	}

	canteen, err := os.OrderRepository.GetCanteenByName(shop.CanteenName)
	if err != nil {
		return nil, err
	}

	dropOff, err := os.OrderRepository.GetDropOffByID(order.DropOffLocationID)
	if err != nil {
		return nil, err
	}

	cLat, _ := strconv.ParseFloat(canteen.Latitude, 64)
	cLon, _ := strconv.ParseFloat(canteen.Longitude, 64)
	dLat, _ := strconv.ParseFloat(dropOff.Latitude, 64)
	dLon, _ := strconv.ParseFloat(dropOff.Longitude, 64)

	distance := calculateDistance(cLat, cLon, dLat, dLon)
	shippingFee := calculateShippingFee(distance)
	var amount float64

	var menuQuantityResp []model.MenuQuantityResponse
	for _, mq := range order.MenuQuantity {
		menu, err := os.OrderRepository.GetMenuByID(mq.MenuID)
		if err != nil {
			return nil, err
		}
		amount += menu.Price * float64(mq.Quantity)
		menuQuantityResp = append(menuQuantityResp, model.MenuQuantityResponse{
			MenuID:   mq.MenuID,
			Quantity: mq.Quantity,
		})
	}

	return &model.GetOrderResponse{
		OrderID:           order.OrderID,
		UserOrderID:       order.UserOrderID,
		UserDeliveryID:    order.UserDeliveryID,
		Status:            order.Status,
		OrderDate:         order.OrderDate,
		DeliveryMethod:    order.DeliveryMethod,
		AppointmentTime:   order.AppointmentTime,
		DropOffLocationID: order.DropOffLocationID,
		MenuQuantity:      menuQuantityResp,
		ShopName:          shop.Name,
		CanteenName:       canteen.CanteenName,
		ShippingFee:       shippingFee,
		Amount:            amount,
	}, nil
}

func (oc *OrderServiceImpl) CreateTransactionLog(log *model.TransactionLog) error {
	targetUserUUID, err := uuid.Parse(log.TargetUserID)
	if err != nil {
		return fmt.Errorf("invalid TargetUserID: %v", err)
	}
	orderUUID, err := uuid.Parse(log.OrderID)
	if err != nil {
		return fmt.Errorf("invalid OrderID: %v", err)
	}

	order, err := oc.OrderRepository.GetOrderByID(orderUUID)
	if err != nil {
		return err
	}
	if order == nil || len(order.MenuQuantity) == 0 {
		return errors.New("order not found or has no menu items")
	}

	firstMenu, err := oc.OrderRepository.GetMenuByID(order.MenuQuantity[0].MenuID)
	if err != nil {
		return err
	}

	var totalMenuPrice float64
	for _, mq := range order.MenuQuantity {
		menu, err := oc.OrderRepository.GetMenuByID(mq.MenuID)
		if err != nil {
			return fmt.Errorf("failed to get menu %v: %w", mq.MenuID, err)
		}
		totalMenuPrice += menu.Price * float64(mq.Quantity)
	}

	shop, err := oc.OrderRepository.GetShopByID(firstMenu.ShopID)
	if err != nil {
		return err
	}

	canteen, err := oc.OrderRepository.GetCanteenByName(shop.CanteenName)
	if err != nil {
		return err
	}

	dropOff, err := oc.OrderRepository.GetDropOffByID(order.DropOffLocationID)
	if err != nil {
		return err
	}

	cLat, _ := strconv.ParseFloat(canteen.Latitude, 64)
	cLon, _ := strconv.ParseFloat(canteen.Longitude, 64)
	dLat, _ := strconv.ParseFloat(dropOff.Latitude, 64)
	dLon, _ := strconv.ParseFloat(dropOff.Longitude, 64)

	distance := calculateDistance(cLat, cLon, dLat, dLon)
	shippingFee := calculateShippingFee(distance)
	totalAmount := totalMenuPrice + shippingFee

	entitiesLog := &entities.TransactionLog{
		TransactionLogID: uuid.New(),
		TargetUserID:     targetUserUUID,
		OrderID:          orderUUID,
		Detail:           log.Detail,
		Amount:           totalAmount,
		TimeStamp:        time.Now(),
	}
	return oc.OrderRepository.CreateTransactionLog(entitiesLog)
}

func (oc *OrderServiceImpl) CreateNotification(notification *model.CreateNotificationRequest) error {
	receiverUUID, err := uuid.Parse(notification.ReceiverID)
	if err != nil {
		return fmt.Errorf("invalid ReceiverID: %v", err)
	}

	orderUUID, err := uuid.Parse(notification.OrderID)
	if err != nil {
		return fmt.Errorf("invalid OrderID: %v", err)
	}

	notificationEntity := &entities.Notification{
		NotificationID: uuid.New(),
		OrderID:        orderUUID,
		ReceiverID:     receiverUUID,
		Topic:          notification.Topic,
		Message:        notification.Message,
		TimeStamp:      time.Now(),
	}

	return oc.OrderRepository.CreateNotification(notificationEntity)
}
