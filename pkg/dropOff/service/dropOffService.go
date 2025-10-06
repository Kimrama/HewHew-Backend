package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/dropOff/model"

	"github.com/google/uuid"
)

type DropOffService interface {
	CreateDropOff(model *model.DropOffRequest) error
	GetAllDropOffs() ([]*entities.DropOff, error)
	GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOff, error)
}
