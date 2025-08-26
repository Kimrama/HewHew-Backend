package entities

import "github.com/google/uuid"

type Menu struct {
	MenuID    uuid.UUID   `gorm:"primaryKey"`
	ShopID    uuid.UUID   `gorm:"not null"`
	Name      string      `gorm:"not null"`
	Price     float64     `gorm:"not null"`
	Status    string      `gorm:"default:'available'"`
	ImageUrl  string      `gorm:"size:512"`
	Favourite []Favourite `gorm:"foreignKey:MenuID"`
}
