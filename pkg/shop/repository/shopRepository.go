package repository

import "hewhew-backend/entities"

type ShopRepository interface {
	CreateCanteen(canteenModel interface{}) error
	EditCanteen(canteenName string,canteenEntity *entities.Canteen) error
	DeleteCanteen(canteenID string) error
}
