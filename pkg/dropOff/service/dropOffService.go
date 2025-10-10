package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/dropOff/model"

	"github.com/google/uuid"
)

type DropOffService interface {
	CreateDropOff(model *model.CreateDropOffRequest) error
	GetAllDropOffs() ([]*entities.DropOffLocation, error)
	GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOffLocation, error)
}
