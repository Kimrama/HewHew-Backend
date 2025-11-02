package repository

import (
	"bytes"
	"fmt"
	"hewhew-backend/config"
	"hewhew-backend/database"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/order/model"
	"hewhew-backend/utils"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const OrderExpirationDuration = 30000000 * time.Minute

type OrderRepositoryImpl struct {
	db             database.Database
	supabaseConfig *config.Supabase
}

func NewOrderRepositoryImpl(db database.Database, supabaseConfig *config.Supabase) OrderRepository {
	return &OrderRepositoryImpl{
		db:             db,
		supabaseConfig: supabaseConfig,
	}
}

func (or *OrderRepositoryImpl) CreateMenuQuantity(menuQuantityEntity *entities.MenuQuantity) error {
	return or.db.Connect().Create(menuQuantityEntity).Error
}

func (or *OrderRepositoryImpl) CreateOrder(orderEntity *entities.Order) error {
	return or.db.Connect().Create(orderEntity).Error
}

func (or *OrderRepositoryImpl) AcceptOrder(acceptOrderModel *model.AcceptOrderRequest) error {
	db := or.db.Connect()
	err := db.Model(&entities.Order{}).
		Where("order_id = ? AND status = ?", acceptOrderModel.OrderID, "waiting").
		Updates(map[string]interface{}{
			"user_delivery_id": acceptOrderModel.DeliveryuserID,
			"status":           "accepted",
		}).Error
	return err
}

func (or *OrderRepositoryImpl) CountActiveOrdersByUser(userID uuid.UUID) (int64, error) {
	db := or.db.Connect()
	var count int64
	err := db.Model(&entities.Order{}).
		Where("user_delivery_id = ? AND status != ?", userID, "delivered").
		Count(&count).Error
	return count, err
}

func (or *OrderRepositoryImpl) ConfirmOrder(orderID uuid.UUID, imageurl string) error {
	db := or.db.Connect()
	err := db.Model(&entities.Order{}).
		Where("order_id = ? AND status = ?", orderID, "accepted").
		Updates(map[string]interface{}{
			"confirmation_image_url": imageurl,
			"status":                 "delivered",
		}).Error
	return err
}

func (or *OrderRepositoryImpl) UploadConfirmImage(orderID uuid.UUID, imageModel *utils.ImageModel) (string, error) {
	customName := orderID.String() + "_" + fmt.Sprintf("%d", time.Now().Unix())

	mimeType := mime.TypeByExtension(imageModel.Ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	url := fmt.Sprintf("%s/storage/v1/object/images/orderConfirmImage/%s", or.supabaseConfig.URL, customName)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(imageModel.Body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", or.supabaseConfig.Key))
	req.Header.Set("Content-Type", mimeType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload image: %s, %s", resp.Status, string(body))
	}
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/images/orderConfirmImage/%s", or.supabaseConfig.URL, customName)
	return publicURL, nil
}

func (or *OrderRepositoryImpl) DeleteOrder(orderID uuid.UUID) error {
	return or.db.Connect().Delete(&entities.Order{}, "order_id = ?", orderID).Error
}

func (or *OrderRepositoryImpl) GetOrdersByShopID(shopID uuid.UUID) ([]*entities.Order, error) {
	db := or.db.Connect()
	var orders []*entities.Order

	err := db.
		Model(&entities.Order{}).
		Joins("JOIN menu_quantities ON menu_quantities.order_id = orders.order_id").
		Joins("JOIN menus ON menus.menu_id = menu_quantities.menu_id").
		Where("menus.shop_id = ? AND orders.status = ?", shopID, "waiting").
		Group("orders.order_id").
		Preload("MenuQuantity").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		if order.Status == "waiting" && time.Since(order.OrderDate) > OrderExpirationDuration {
			order.Status = "expired"
			db.Model(&order).Update("status", "expired")
		}
	}

	return orders, nil
}

func (or *OrderRepositoryImpl) GetOrdersByCanteens(canteens []string) ([]entities.Order, error) {
	db := or.db.Connect()
	var orders []entities.Order

	orderClause := "CASE"
	for i, c := range canteens {
		orderClause += fmt.Sprintf(" WHEN s.canteen_name='%s' THEN %d", c, i)
	}
	orderClause += " END"

	err := db.
		Joins("JOIN menu_quantities mq ON mq.order_id = orders.order_id").
		Joins("JOIN menus m ON m.menu_id = mq.menu_id").
		Joins("JOIN shops s ON s.shop_id = m.shop_id").
		Where("s.canteen_name IN ?", canteens).
		Where("orders.status = ?", "waiting").
		Order(orderClause).
		Preload("MenuQuantity").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		if order.Status == "waiting" && time.Since(order.OrderDate) > OrderExpirationDuration {
			order.Status = "expired"
			db.Model(&order).Update("status", "expired")
		}
	}

	return orders, nil
}

func (or *OrderRepositoryImpl) GetOrdersByUserID(userID uuid.UUID) ([]*entities.Order, error) {
	db := or.db.Connect()
	var orders []*entities.Order

	err := db.Where("user_order_id = ?", userID).
		Preload("MenuQuantity").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		if order.Status == "waiting" && time.Since(order.OrderDate) > OrderExpirationDuration {
			order.Status = "expired"
			db.Model(&order).Update("status", "expired")
		}
	}

	return orders, nil
}

func (or *OrderRepositoryImpl) GetOrderByDeliveryUserID(userID uuid.UUID) ([]*entities.Order, error) {
	db := or.db.Connect()
	var orders []*entities.Order

	err := db.Where("user_delivery_id = ?", userID).
		Preload("MenuQuantity").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		if order.Status == "waiting" && time.Since(order.OrderDate) > OrderExpirationDuration {
			order.Status = "expired"
			db.Model(&order).Update("status", "expired")
		}
	}

	return orders, nil
}

func (or *OrderRepositoryImpl) GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error) {
	var admin entities.ShopAdmin
	db := or.db.Connect()

	if err := db.Select("shop_id").First(&admin, "admin_id = ?", adminID).Error; err != nil {
		return nil, err
	}
	var shop entities.Shop
	if err := db.First(&shop, "shop_id = ?", admin.ShopID).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

func (or *OrderRepositoryImpl) GetAvailableOrders() ([]*entities.Order, error) {
	db := or.db.Connect()
	var orders []*entities.Order

	if err := db.Preload("MenuQuantity").
		Where("status = ?", "waiting").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	for _, order := range orders {
		if time.Since(order.OrderDate) > OrderExpirationDuration {
			order.Status = "expired"
			db.Model(&order).Update("status", "expired")
		}
	}

	return orders, nil
}

func (or *OrderRepositoryImpl) CreateReview(reviewEntity *entities.Review) error {
	err := or.db.Connect().Create(reviewEntity).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *OrderRepositoryImpl) GetReviewsByTargetUserID(userID uuid.UUID) ([]*entities.Review, error) {
	db := or.db.Connect()
	var reviews []*entities.Review
	err := db.Where("user_target_id = ?", userID).
		Order("time_stamp DESC").
		Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (or *OrderRepositoryImpl) GetReviewsByReviewerUserID(userID uuid.UUID) ([]*entities.Review, error) {
	db := or.db.Connect()
	var reviews []*entities.Review
	err := db.Where("user_reviewer_id = ?", userID).
		Order("time_stamp DESC").
		Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (or *OrderRepositoryImpl) GetReviewByID(reviewID uuid.UUID) (*entities.Review, error) {
	db := or.db.Connect()
	var review entities.Review
	err := db.Where("review_id = ?", reviewID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (or *OrderRepositoryImpl) GetOrderByID(orderID uuid.UUID) (*entities.Order, error) {
	db := or.db.Connect()
	var order entities.Order
	err := db.Preload("MenuQuantity").
		First(&order, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}

	if order.Status == "waiting" && time.Since(order.OrderDate) > OrderExpirationDuration {
		order.Status = "expired"
		db.Model(&order).Update("status", "expired")
	}

	return &order, nil
}

func (or *OrderRepositoryImpl) GetMenuByID(menuID uuid.UUID) (*entities.Menu, error) {
	db := or.db.Connect()
	var menu entities.Menu
	err := db.First(&menu, "menu_id = ?", menuID).Error
	return &menu, err
}

func (or *OrderRepositoryImpl) GetShopByID(shopID uuid.UUID) (*entities.Shop, error) {
	db := or.db.Connect()
	var shop entities.Shop
	err := db.First(&shop, "shop_id = ?", shopID).Error
	return &shop, err
}

func (or *OrderRepositoryImpl) GetCanteenByName(name string) (*entities.Canteen, error) {
	db := or.db.Connect()
	var canteen entities.Canteen
	err := db.First(&canteen, "canteen_name = ?", name).Error
	return &canteen, err
}

func (or *OrderRepositoryImpl) GetDropOffByID(id uuid.UUID) (*entities.DropOffLocation, error) {
	db := or.db.Connect()
	var dropOff entities.DropOffLocation
	err := db.First(&dropOff, "drop_off_location_id = ?", id).Error
	return &dropOff, err
}

func (or *OrderRepositoryImpl) GetAllCanteens() ([]*entities.Canteen, error) {
	db := or.db.Connect()
	var canteens []*entities.Canteen
	err := db.Find(&canteens).Error
	return canteens, err
}

func (or *OrderRepositoryImpl) CreateNotification(notification *entities.Notification) error {
	return or.db.Connect().Create(notification).Error
}

func (or *OrderRepositoryImpl) CreateNotificationDriver(notification *entities.Notification) error {
	return or.db.Connect().Create(notification).Error
}

func (or *OrderRepositoryImpl) CreateTransactionLog(log *entities.TransactionLog) error {
	return or.db.Connect().Create(log).Error
}

func (or *OrderRepositoryImpl) CheckReviewExists(orderID, userReviewerID uuid.UUID) (bool, error) {
	var count int64
	err := or.db.Connect().
		Model(&entities.Review{}).
		Where("order_id = ? AND user_reviewer_id = ?", orderID, userReviewerID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (or *OrderRepositoryImpl) GetNotificationByUserID(userID uuid.UUID) ([]*entities.Notification, error) {
	db := or.db.Connect()
	var notifications []*entities.Notification
	err := db.Where("receiver_id = ?", userID).
		Order("time_stamp DESC").
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (or *OrderRepositoryImpl) GetUserByID(userID uuid.UUID) (*entities.User, error) {
	var user entities.User
	if err := or.db.Connect().First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (or *OrderRepositoryImpl) UpdateWalletBalance(userID uuid.UUID, newBalance float64) error {
	return or.db.Connect().
		Model(&entities.User{}).
		Where("user_id = ?", userID).
		Update("wallet", newBalance).Error
}