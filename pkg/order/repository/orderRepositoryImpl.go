package repository

import (
	"hewhew-backend/database"
	"hewhew-backend/entities"

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

func (or *OrderRepositoryImpl) CreateOrder(orderEntity *entities.Order) error {
	return or.db.Connect().Create(orderEntity).Error
}

func (or *OrderRepositoryImpl) CreateMenuQuantity(menuQuantityEntity *entities.MenuQuantity) error {
	return or.db.Connect().Create(menuQuantityEntity).Error
}

func (or *OrderRepositoryImpl) GetMenuByID(menuID uuid.UUID) (*entities.Menu, error) {
	var menu entities.Menu
	err := or.db.Connect().Where("menu_id = ?", menuID).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (or *OrderRepositoryImpl) GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOff, error) {
	var dropOff entities.DropOff
	err := or.db.Connect().Where("drop_off_id = ?", dropOffID).First(&dropOff).Error
	if err != nil {
		return nil, err
	}
	return &dropOff, nil
}
