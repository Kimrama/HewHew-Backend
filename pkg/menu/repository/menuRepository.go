package repository

type MenuRepository interface {
	CreateMenu(menuModel interface{}) error
}
