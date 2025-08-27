package entities

import "github.com/google/uuid"

type Favourite struct {
	FavouriteID uuid.UUID `gorm:"primaryKey"`
	UserID      uuid.UUID `gorm:"not null;index:idx_user_menu,unique"`
	MenuID      uuid.UUID `gorm:"not null;index:idx_user_menu,unique"`
}
