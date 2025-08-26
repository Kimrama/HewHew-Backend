package entities

import (
	"time"

	"github.com/google/uuid"
)

type NotificationTransaction struct {
	NotificationTransactionID uuid.UUID `gorm:"primaryKey"`
	TimeStamp                 time.Time `gorm:"autoCreateTime:milli"`
}
