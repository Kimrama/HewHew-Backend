package entities

import "github.com/google/uuid"

type Contact struct {
	ContactID   uuid.UUID `gorm:"primaryKey"`
	UserID      uuid.UUID `gorm:"not null"`
	ContactType string    `gorm:"not null"`
	Detail      string    `gorm:"not null"`
}
