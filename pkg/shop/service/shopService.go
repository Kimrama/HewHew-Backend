package service

type ShopService interface {
	CreateCanteen(canteenModel interface{}) error
	EditCanteen(canteenModel interface{}) error
	DeleteCanteen(canteenID string) error
	GetCanteens() (interface{}, error)
}
