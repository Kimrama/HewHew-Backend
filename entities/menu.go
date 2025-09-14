package entities

import "github.com/google/uuid"

type Menu struct {
	MenuID    uuid.UUID   `gorm:"primaryKey"`	
	ShopID    uuid.UUID   `gorm:"not null;index"`
	Name      string      `gorm:"not null"`
	Detail    string      `gorm:"type:text"`
	Price     float64     `gorm:"not null"`
	Status    string      `gorm:"default:'available'"`
	ImageURL  string      `gorm:"size:512"`
	Favourite []Favourite `gorm:"foreignKey:MenuID"`
	MenuQuantity []MenuQuantity `gorm:"foreignKey:MenuID"`
	Tag1ID     *uuid.UUID   `gorm:"index"`
	Tag2ID     *uuid.UUID   `gorm:"index"`
}
