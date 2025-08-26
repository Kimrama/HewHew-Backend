package entities

import "github.com/google/uuid"

type TopUp struct {
	TopUpID uuid.UUID `gorm:"primaryKey"`
	UserID  uuid.UUID `gorm:"not null"`
	Amount  float64   `gorm:"not null"`
}
