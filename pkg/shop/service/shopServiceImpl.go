package service

import (
	"errors"
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/shop/model"
	"hewhew-backend/pkg/shop/repository"
	"hewhew-backend/utils"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ShopServiceImpl struct {
	ShopRepository repository.ShopRepository
}

func NewShopServiceImpl(ShopRepository repository.ShopRepository) ShopService {
	return &ShopServiceImpl{
		ShopRepository: ShopRepository,
	}
}
func (s *ShopServiceImpl) CreateCanteen(canteenModel interface{}) error {
	return s.ShopRepository.CreateCanteen(canteenModel)
}

func (s *ShopServiceImpl) EditCanteen(canteenName string, canteenEntity *entities.Canteen) error {
	if canteenName == "" {
		return fmt.Errorf("canteen name is required")
	}
	return s.ShopRepository.EditCanteen(canteenName, canteenEntity)
}

func (s *ShopServiceImpl) DeleteCanteen(canteenName string) error {
	if canteenName == "" {
		return fmt.Errorf("canteen name is required")
	}
	return s.ShopRepository.DeleteCanteen(canteenName)
}

func (s *ShopServiceImpl) GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error) {
	return s.ShopRepository.GetShopByAdminID(adminID)
}

func (s *ShopServiceImpl) ChangeState(body model.ChangeState, shopID uuid.UUID) error {

	var state bool
	switch body.State {
	case "open":
		state = true
	case "close":
		state = false
	default:
		return fmt.Errorf("invalid state value: %s", body.State)
	}

	return s.ShopRepository.ChangeState(state, shopID)
}

func (s *ShopServiceImpl) EditShopImage(shopID uuid.UUID, imageModel *utils.ImageModel) error {
	err := s.ShopRepository.EditShopImage(shopID, imageModel)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShopServiceImpl) EditShop(body model.EditShopRequest, shop uuid.UUID) error {
	if body.ShopName == "" && body.ShopCanteenName == "" {
		return errors.New("no fields to update")
	}

	shopEntity := &entities.Shop{
		Name:        body.ShopName,
		CanteenName: body.ShopCanteenName,
		Address:     "Null",
	}
	fmt.Println("Service - EditShop: ", shopEntity, shop)

	return s.ShopRepository.EditShop(*shopEntity, shop)
}

func (s *ShopServiceImpl) CreateTag(ShopID string, body *model.TagCreateRequest) (*entities.Tag, error) {
	if body.Topic == "" {
		return nil, errors.New("tag topic is required")
	}

	shopUUID, err := uuid.Parse(ShopID)
	if err != nil {
		return nil, fmt.Errorf("invalid ShopID: %v", err)
	}

	tagEntity := &entities.Tag{
		Topic:  body.Topic,
		ShopID: shopUUID,
		TagID:  uuid.New(),
	}

	return s.ShopRepository.CreateTag(tagEntity)
}

func (s *ShopServiceImpl) GetShopAdminByUsername(username string) (*entities.ShopAdmin, error) {
	return s.ShopRepository.GetShopAdminByUsername(username)
}

func (s *ShopServiceImpl) GetAllCanteens() ([]entities.Canteen, error) {
	return s.ShopRepository.GetAllCanteens()
}

func (s *ShopServiceImpl) GetCanteenByName(canteenName string) (*entities.Canteen, error) {
	return s.ShopRepository.GetCanteenByName(canteenName)
}

func (s *ShopServiceImpl) GetAllShops() ([]entities.Shop, error) {
	return s.ShopRepository.GetAllShops()
}

func (s *ShopServiceImpl) GetShopByID(shopID uuid.UUID) (*entities.Shop, error) {
	return s.ShopRepository.GetShopByID(shopID)
}

func (s *ShopServiceImpl) GetTagsByShopIDAndTopic(shopID string, topic string) ([]entities.Tag, error) {
	return s.ShopRepository.GetTagsByShopIDAndTopic(shopID, topic)
}

func (s *ShopServiceImpl) EditTag(tagID string, topic string) error {
	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		return fmt.Errorf("invalid TagID: %v", err)
	}

	tagEntity := &entities.Tag{
		TagID: tagUUID,
		Topic: topic,
	}

	return s.ShopRepository.EditTag(tagEntity)
}

func (s *ShopServiceImpl) GetAllTags(shopID uuid.UUID) ([]entities.Tag, error) {
	return s.ShopRepository.GetAllTags(shopID)
}
func (s *ShopServiceImpl) DeleteTag(tagID string) error {
	err := s.ShopRepository.DeleteTag(tagID)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %v", err)
	}
	return nil
}

func (s *ShopServiceImpl) GetAllMenus(shopID uuid.UUID) ([]*entities.Menu, error) {
	return s.ShopRepository.GetAllMenus(shopID)
}

func (s *ShopServiceImpl) CreateTransactionLog(log *model.TransactionLog) error {
	targetUserUUID, err := uuid.Parse(log.TargetUserID)
	if err != nil {
		return fmt.Errorf("invalid TargetUserID: %v", err)
	}
	orderUUID, err := uuid.Parse(log.OrderID)
	if err != nil {
		return fmt.Errorf("invalid OrderID: %v", err)
	}

	order, err := s.ShopRepository.GetOrderByID(orderUUID)
	if err != nil {
		return err
	}
	if order == nil || len(order.MenuQuantity) == 0 {
		return errors.New("order not found or has no menu items")
	}

	firstMenu, err := s.ShopRepository.GetMenuByID(order.MenuQuantity[0].MenuID)
	if err != nil {
		return err
	}

	var totalMenuPrice float64
	for _, mq := range order.MenuQuantity {
		menu, err := s.ShopRepository.GetMenuByID(mq.MenuID)
		if err != nil {
			return fmt.Errorf("failed to get menu %v: %w", mq.MenuID, err)
		}
		totalMenuPrice += menu.Price * float64(mq.Quantity)
	}

	shop, err := s.ShopRepository.GetShopByID(firstMenu.ShopID)
	if err != nil {
		return err
	}

	canteen, err := s.ShopRepository.GetCanteenByName(shop.CanteenName)
	if err != nil {
		return err
	}

	dropOff, err := s.ShopRepository.GetDropOffByID(order.DropOffLocationID)
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
	return s.ShopRepository.CreateTransactionLog(entitiesLog)
}

func (s *ShopServiceImpl) CreateNotification(notification *model.CreateNotificationRequest) error {
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

	return s.ShopRepository.CreateNotification(notificationEntity)
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
