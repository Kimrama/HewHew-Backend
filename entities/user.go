package entities

type User struct {
	UserID          string  `gorm:"primaryKey"`
	UserName        string  `gorm:"not null"`
	Password        string  `gorm:"not null"`
	FName           string  `gorm:"not null"`
	LName           string  `gorm:"not null"`
	Faculty         string  `gorm:"not null"`
	ProfileImageUrl string  `gorm:"size:512"`
	Gender          string  `gorm:"default:'undefined'"`
	Wallet          float64 `gorm:"default:0"`
}
