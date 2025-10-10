package repository

import (
	"hewhew-backend/entities"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type DropOffRepository interface {
	CreateDropOff(dropOff *entities.DropOffLocation) error
	UploadDropOffImage(DropOffLocationID uuid.UUID, imageModel *utils.ImageModel) (string, error)
	GetAllDropOffs() ([]*entities.DropOffLocation, error)
	GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOffLocation, error)
}
