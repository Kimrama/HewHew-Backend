package entities

import "github.com/google/uuid"

type Contract struct {
	ContractID   uuid.UUID `gorm:"primaryKey"`
	UserID       uuid.UUID `gorm:"not null"`
	ContractType string    `gorm:"not null"`
	Detail       string    `gorm:"not null"`
}
