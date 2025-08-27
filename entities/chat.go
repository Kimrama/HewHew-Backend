package entities

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ChatID    uuid.UUID `gorm:"primaryKey"`
	Sender    uuid.UUID `gorm:"not null"`
	Message   string    `gorm:"not null"`
	OrderID   uuid.UUID `gorm:"not null"`
	TimeStamp time.Time `gorm:"autoCreateTime:milli"`
}
