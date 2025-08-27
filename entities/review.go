package entities

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ReviewID       uuid.UUID `gorm:"primaryKey"`
	UserReviewerID uuid.UUID `gorm:"not null"`
	UserTargetID   uuid.UUID `gorm:"not null"`
	OrderID        uuid.UUID `gorm:"not null"`
	Rating         int       `gorm:"not null"`
	Comment        string    `gorm:"size:512"`
	TimeStamp      time.Time `gorm:"autoCreateTime:milli"`
}
