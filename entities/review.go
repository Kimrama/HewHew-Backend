package entities

import "github.com/google/uuid"

type Review struct {
	ReviewID uuid.UUID `gorm:"primaryKey"`
	OrderID  uuid.UUID `gorm:"not null"`
	Rating   int       `gorm:"not null"`
	Comment  string    `gorm:"size:512"`
}
