package entities

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID              uuid.UUID      `gorm:"primaryKey"`
	UserOrderID          uuid.UUID      `gorm:"not null"`
	UserDeliveryID       uuid.UUID      `gorm:"default:'NULL'"`
	Status               string         `gorm:"index;default:'Un Paid'"`
	OrderDate            time.Time      `gorm:"autoCreateTime:milli"`
	DeliveryMethod       string         `gorm:"not null"`
	ConfirmationImageURL string         `gorm:"size:512"`
	AppointmentTime      time.Time      `gorm:"not null"`
	DropOffLocation      string         `gorm:"size:256"`
	MenuQuantity         []MenuQuantity `gorm:"foreignKey:OrderID"`
	TransactionLog       TransactionLog `gorm:"foreignKey:OrderID"`
	Notifications        []Notification `gorm:"foreignKey:OrderID"`
	Chats                []Chat         `gorm:"foreignKey:OrderID"`
}
