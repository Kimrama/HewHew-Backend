package service

import "hewhew-backend/entities"

type ShopService interface {
	CreateCanteen(canteenModel interface{}) error
	EditCanteen(canteenName string,canteenEntity *entities.Canteen) error
	DeleteCanteen(canteenName string) error
}
