package service

type MenuService interface {
	CreateMenu(menuModel interface{}) error
}
