package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type MenuRepository interface {
	CreateMenu(menuEntity *entities.Menu) error
	UploadMenuImage(menuID uuid.UUID, imageModel *utils.ImageModel) (string, error)
}
