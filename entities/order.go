package entities

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID              uuid.UUID      `gorm:"primaryKey"`
	UserOrderID          uuid.UUID      `gorm:"not null"`
	UserDeliveryID       *uuid.UUID     `gorm:"default:null"`
	Status               string         `gorm:"index;default:'waiting'"`
	OrderDate            time.Time      `gorm:"autoCreateTime:milli"`
	DeliveryMethod       string         `gorm:"not null"`
	ConfirmationImageURL string         `gorm:"size:512"`
	AppointmentTime      time.Time      `gorm:"not null"`
	DropOffLocationID    uuid.UUID      `gorm:"not null"`
	MenuQuantity         []MenuQuantity `gorm:"foreignKey:OrderID"`
	TransactionLog       TransactionLog `gorm:"foreignKey:OrderID"`
	Notifications        []Notification `gorm:"foreignKey:OrderID"`
	Chats                []Chat         `gorm:"foreignKey:OrderID"`
}
