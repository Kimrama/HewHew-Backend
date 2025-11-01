package entities

import "github.com/google/uuid"

type Shop struct {
	ShopID      uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Address     string    `gorm:"not null"`
	ImageURL    string    `gorm:"size:512"`
	Menus       []Menu    `gorm:"foreignKey:ShopID"`
	Tags        []Tag     `gorm:"foreignKey:ShopID"`
	CanteenName string    `gorm:"not null;index:idx_canteen_shop"`
	State       bool      `gorm:"not null;default:true"`
}
