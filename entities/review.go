package entities

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ReviewID       uuid.UUID `gorm:"primaryKey"`
	UserReviewerID uuid.UUID `gorm:"not null;index:idx_reviewer_target"`
	UserTargetID   uuid.UUID `gorm:"not null;index:idx_reviewer_target"`
	OrderID        uuid.UUID `gorm:"not null"`
	Rating         int       `gorm:"not null"`
	Comment        string    `gorm:"size:512"`
	TimeStamp      time.Time `gorm:"autoCreateTime:milli"`
}
