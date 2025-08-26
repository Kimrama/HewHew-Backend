package entities

import "github.com/google/uuid"

type MenuQuantity struct {
	MenuQuantityID uuid.UUID `gorm:"primaryKey"`
	MenuID         uuid.UUID `gorm:"not null"`
	OrderID        uuid.UUID `gorm:"not null"`
	Quantity       int       `gorm:"not null"`
}
