package repository

import (
	"bytes"
	"database/sql"
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

func (or *OrderRepositoryImpl) GetUserAverageRating(userID uuid.UUID) (float64, error) {
	db := or.db.Connect()
	var avg sql.NullFloat64

	err := db.Table("reviews").
		Select("AVG(reviews.rating)").
		Joins("JOIN orders ON reviews.order_id = orders.order_id").
		Where("reviews.user_target_id = ? AND orders.user_delivery_id = ?", userID, userID).
		Scan(&avg).Error

	if err != nil {
		return 0, err
	}
	if !avg.Valid {
		return 0, nil
	}
	return avg.Float64, nil
}

func (or *OrderRepositoryImpl) CountActiveOrdersByUser(userID uuid.UUID) (int64, error) {
	db := or.db.Connect()
	var count int64
	err := db.Model(&entities.Order{}).
		Where("user_delivery_id = ? AND status != ?", userID, "finished").
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

func (or *OrderRepositoryImpl) GetOrderByID(orderID uuid.UUID) (*entities.Order, error) {
	var order entities.Order
	err := or.db.Connect().Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *OrderRepositoryImpl) GetOrdersByShopID(shopID uuid.UUID) ([]*entities.Order, error) {
	var orders []*entities.Order
	db := or.db.Connect()
	err := db.
		Model(&entities.Order{}).
		Joins("JOIN menu_quantities ON menu_quantities.order_id = orders.order_id").
		Joins("JOIN menus ON menus.menu_id = menu_quantities.menu_id").
		Where("menus.shop_id = ? AND orders.status = ?", shopID, "waiting").
		Group("orders.order_id").
		Preload("MenuQuantity").
		Find(&orders).Error

	return orders, err
}

func (or *OrderRepositoryImpl) GetOrdersByUserID(userID uuid.UUID) ([]*entities.Order, error) {
	var orders []*entities.Order
	db := or.db.Connect()
	err := db.Where("user_order_id = ?", userID).Preload("MenuQuantity").Find(&orders).Error
	return orders, err
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
	var orders []*entities.Order
	db := or.db.Connect()
	err := db.Where("status = ?", "waiting").Preload("MenuQuantity").Find(&orders).Error
	return orders, err
}

func (or *OrderRepositoryImpl) GetMenuByID(menuID uuid.UUID) (*entities.Menu, error) {
	var menu entities.Menu
	err := or.db.Connect().Where("menu_id = ?", menuID).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (or *OrderRepositoryImpl) GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOffLocation, error) {
	var dropOff entities.DropOffLocation
	err := or.db.Connect().Where("drop_off_location_id = ?", dropOffID).First(&dropOff).Error
	if err != nil {
		return nil, err
	}
	return &dropOff, nil
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

func (or *OrderRepositoryImpl) GetReviewByID(reviewID uuid.UUID) (*entities.Review, error) {
	db := or.db.Connect()
	var review entities.Review
	err := db.Where("review_id = ?", reviewID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}
