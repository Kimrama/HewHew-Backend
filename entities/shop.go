package entities

import "github.com/google/uuid"

type Shop struct {
	ShopID   uuid.UUID `gorm:"primaryKey"`
	Name     string    `gorm:"not null"`
	Address  string    `gorm:"not null"`
	ImageUrl string    `gorm:"size:512"`
	Menus    []Menu    `gorm:"foreignKey:ShopID"`
}
