package entities

import "github.com/google/uuid"

type DropOff struct {
	DropOffID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Latitude  string    `gorm:"not null"`
	Longitude string    `gorm:"not null"`
	Orders    []Order   `gorm:"foreignKey:DropOffID"`
}
