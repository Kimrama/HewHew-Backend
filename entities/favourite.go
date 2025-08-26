package entities

import "github.com/google/uuid"

type Favourite struct {
	FavouriteID uuid.UUID `gorm:"primaryKey"`
	UserID      uuid.UUID `gorm:"not null"`
	MenuID      uuid.UUID `gorm:"not null"`
}
