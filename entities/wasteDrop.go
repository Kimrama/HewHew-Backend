package entities

type WasteDrop struct {
	Latitude  float64 `gorm:"primaryKey"`
	Longitude float64 `gorm:"primaryKey"`
}
