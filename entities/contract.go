package entities

type Contract struct {
	ContractID   string `gorm:"primaryKey"`
	UserID       string `gorm:"not null"`
	ContractType string `gorm:"not null"`
	Detail       string `gorm:"not null"`
}
