package repository

import (
	"hewhew-backend/database"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/order/model"

	"github.com/google/uuid"
)

type OrderRepositoryImpl struct {
	db database.Database
}

func NewOrderRepositoryImpl(db database.Database) OrderRepository {
	return &OrderRepositoryImpl{
		db: db,
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
