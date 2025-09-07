package entities

type Canteen struct {
	CanteenName string `gorm:"primaryKey"`
	Latitude    string `gorm:"not null"`
	Longitude   string `gorm:"not null"`
	Shops       []Shop `gorm:"foreignKey:CanteenName"`
}
