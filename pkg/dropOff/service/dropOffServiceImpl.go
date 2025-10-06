package service

import (
	"hewhew-backend/entities"
	"hewhew-backend/pkg/dropOff/model"
	"hewhew-backend/pkg/dropOff/repository"

	"github.com/google/uuid"
)

type DropOffServiceImpl struct {
	DropOffRepository repository.DropOffRepository
}

func NewDropOffServiceImpl(DropOffRepository repository.DropOffRepository) DropOffService {
	return &DropOffServiceImpl{
		DropOffRepository: DropOffRepository,
	}
}

func (ds *DropOffServiceImpl) CreateDropOff(model *model.DropOffRequest) error {

	do := &entities.DropOff{
		DropOffID: uuid.New(),
		Latitude:  model.Latitude,
		Longitude: model.Longitude,
	}
	return ds.DropOffRepository.CreateDropOff(do)
}

func (ds *DropOffServiceImpl) GetAllDropOffs() ([]*entities.DropOff, error) {
	return ds.DropOffRepository.GetAllDropOffs()
}

func (ds *DropOffServiceImpl) GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOff, error) {
	return ds.DropOffRepository.GetDropOffByID(dropOffID)
}
