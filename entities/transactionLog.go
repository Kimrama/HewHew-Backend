package entities

import (
	"time"

	"github.com/google/uuid"
)

type TransactionLog struct {
	TransactionLogID uuid.UUID `gorm:"primaryKey"`
	TargetUserID     uuid.UUID `gorm:"not null"`
	OrderID          uuid.UUID `gorm:"uniqueIndex"`
	TimeStamp        time.Time `gorm:"autoCreateTime:milli"`
	Detail           string    `gorm:"size:256;not null"`
	Amount           float64   `gorm:"not null"`
}
