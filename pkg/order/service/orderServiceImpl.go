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
	for _, item := range orderModel.Menu {
		menu, err := os.OrderRepository.GetMenuByID(item.MenuID)
		if err != nil {
			return fmt.Errorf("menu with ID %s not found", item.MenuID)
		}
		if menu.Status != "available" {
			return fmt.Errorf("menu %s is not available", menu.Name)
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

func (os *OrderServiceImpl) GetOrdersByUserID(userID uuid.UUID) ([]*entities.Order, error) {
	return os.OrderRepository.GetOrdersByUserID(userID)
}

func (os *OrderServiceImpl) GetOrdersByShopID(userID string) ([]*entities.Order, error) {
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
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	return orders, nil
}

func (os *OrderServiceImpl) GetAvailableOrders() ([]*entities.Order, error) {
	orders, err := os.OrderRepository.GetAvailableOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	return orders, nil
}

func (os *OrderServiceImpl) GetOrderByID(orderID uuid.UUID) (*model.OrderResponse, error) {
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

	var menuQuantityResp []model.MenuQuantityResponse
	for _, mq := range order.MenuQuantity {
		menuQuantityResp = append(menuQuantityResp, model.MenuQuantityResponse{
			MenuID:   mq.MenuID,
			Quantity: mq.Quantity,
		})
	}

	return &model.OrderResponse{
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

func (os *OrderServiceImpl) GetReviewsByTargetUserID(userID uuid.UUID) ([]*entities.Review, error) {
	return os.OrderRepository.GetReviewsByTargetUserID(userID)
}

func (os *OrderServiceImpl) GetReviewByID(reviewID uuid.UUID) (*entities.Review, error) {
	return os.OrderRepository.GetReviewByID(reviewID)
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
