package entities

import "github.com/google/uuid"

type ShopAdmin struct {
	AdminID  uuid.UUID `gorm:"primaryKey"`
	Username string    `gorm:"not null;unique"`
	Password string    `gorm:"not null"`
	FName    string    `gorm:"not null"`
	LName    string    `gorm:"not null"`
	Shop     Shop      `gorm:"foreignKey:ShopID"`
}
