package service

import (
	"hewhew-backend/pkg/menu/model"

	"github.com/google/uuid"
)

type MenuService interface {
	CreateMenu(menuModel *model.CreateMenuRequest,shopID uuid.UUID) error
}
