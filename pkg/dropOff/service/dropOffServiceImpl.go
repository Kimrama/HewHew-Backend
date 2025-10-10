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

func (ds *DropOffServiceImpl) CreateDropOff(model *model.CreateDropOffRequest) error {
	DropOffLocationID := uuid.New()
	imageUrl := ""
	if model.Image != nil {
		var err error
		imageUrl, err = ds.DropOffRepository.UploadDropOffImage(DropOffLocationID, model.Image)
		if err != nil {
			return err
		}
	}
	dropoffEntity := &entities.DropOffLocation{
		DropOffLocationID: DropOffLocationID,
		Latitude:          model.Latitude,
		Longitude:         model.Longitude,
		Name:              model.Name,
		Detail:            model.Detail,
		ImageURL:          imageUrl,
	}
	if err := ds.DropOffRepository.CreateDropOff(dropoffEntity); err != nil {
		return err
	}

	return nil
}

func (ds *DropOffServiceImpl) GetAllDropOffs() ([]*entities.DropOffLocation, error) {
	return ds.DropOffRepository.GetAllDropOffs()
}

func (ds *DropOffServiceImpl) GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOffLocation, error) {
	return ds.DropOffRepository.GetDropOffByID(dropOffID)
}
