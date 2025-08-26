package entities

import (
	"github.com/google/uuid"
)

type Chat struct {
	ChatID  uuid.UUID `gorm:"primaryKey"`
	Sender  uuid.UUID `gorm:"not null"`
	Message string    `gorm:"not null"`
}
