package entities

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	NotificationID uuid.UUID `gorm:"primaryKey"`
	OrderID        uuid.UUID `gorm:"not null"`
	ReceiverID     uuid.UUID `gorm:"not null"`
	Topic          string    `gorm:"size:255;not null"`
	Message        string    `gorm:"size:1024"`
	TimeStamp      time.Time `gorm:"autoCreateTime:milli"`
}
