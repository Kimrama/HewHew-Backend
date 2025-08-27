package entities

type WasteDrop struct {
	Latitude  string `gorm:"primaryKey"`
	Longitude string `gorm:"primaryKey"`
}
