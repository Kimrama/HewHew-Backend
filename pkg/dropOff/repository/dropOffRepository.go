package repository

import (
	"hewhew-backend/entities"

	"github.com/google/uuid"
)

type DropOffRepository interface {
	CreateDropOff(dropOff *entities.DropOff) error
	GetAllDropOffs() ([]*entities.DropOff, error)
	GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOff, error)
}
