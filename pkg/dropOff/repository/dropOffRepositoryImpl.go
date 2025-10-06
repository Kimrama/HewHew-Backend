package repository

import (
	"hewhew-backend/database"
	"hewhew-backend/entities"

	"github.com/google/uuid"
)

type DropOffRepositoryImpl struct {
	db database.Database
}

func NewDropOffRepositoryImpl(db database.Database) DropOffRepository {
	return &DropOffRepositoryImpl{
		db: db,
	}
}

func (dr *DropOffRepositoryImpl) CreateDropOff(dropOff *entities.DropOff) error {
	return dr.db.Connect().Create(dropOff).Error
}

func (dr *DropOffRepositoryImpl) GetAllDropOffs() ([]*entities.DropOff, error) {
	var dropOffs []*entities.DropOff
	err := dr.db.Connect().Find(&dropOffs).Error
	if err != nil {
		return nil, err
	}
	return dropOffs, nil
}
func (dr *DropOffRepositoryImpl) GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOff, error) {
	var dropOff entities.DropOff
	err := dr.db.Connect().Where("drop_off_id = ?", dropOffID).First(&dropOff).Error
	if err != nil {
		return nil, err
	}
	return &dropOff, nil
}
