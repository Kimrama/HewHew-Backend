package entities

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID              uuid.UUID `gorm:"primaryKey"`
	UserOrderID          uuid.UUID `gorm:"not null"`
	UserDeliveryID       uuid.UUID `gorm:"default:'None'"`
	Status               string    `gorm:"default:'Pending'"`
	OrderDate            time.Time `gorm:"autoCreateTime:milli"`
	DeliveryMethod       string    `gorm:"not null"`
	ConfirmationImageUrl string    `gorm:"size:512"`
	AppointmentTime      time.Time `gorm:"not null"`
	DropOffLocation      string    `gorm:"size:256"`
}
