package entities

import "github.com/google/uuid"

type User struct {
	UserID          uuid.UUID        `gorm:"primaryKey"`
	UserName        string           `gorm:"not null"`
	Password        string           `gorm:"not null"`
	FName           string           `gorm:"not null"`
	LName           string           `gorm:"not null"`
	Faculty         string           `gorm:"not null"`
	ProfileImageURL string           `gorm:"size:512"`
	Gender          string           `gorm:"default:'undefined'"`
	Contacts        []Contact        `gorm:"foreignKey:UserID"`
	TopUps          []TopUp          `gorm:"foreignKey:UserID"`
	Orders          []Order          `gorm:"foreignKey:UserOrderID"`
	TransactionLogs []TransactionLog `gorm:"foreignKey:TargetUserID"`
	Notifications   []Notification   `gorm:"foreignKey:ReceiverID"`
	Chats           []Chat           `gorm:"foreignKey:SenderID"`
	Favourites      []Favourite      `gorm:"foreignKey:UserID"`
	Reviews         []Review         `gorm:"foreignKey:UserReviewerID"`
	Wallet          float64          `gorm:"default:0"`
}
