package entities

import "github.com/google/uuid"

type DropOffLocation struct {
	DropOffLocationID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Latitude          string    `gorm:"not null"`
	Longitude         string    `gorm:"not null"`
	Name              string
	Detail            string
	ImageURL          string
	Orders            []Order `gorm:"foreignKey:DropOffLocationID"`
}
