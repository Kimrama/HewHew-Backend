package entities

import "github.com/google/uuid"

type Tag struct {
	TagID  		uuid.UUID `gorm:"primaryKey"`
    Topic   	string    `gorm:"not null;index:idx_shop_topic,unique"`
    ShopID  	uuid.UUID `gorm:"not null;index:idx_shop_topic,unique"`
	Menus1	   	[]Menu `gorm:"foreignKey:Tag1ID"`
	Menus2	   	[]Menu `gorm:"foreignKey:Tag2ID"`
	
}

func (t *Tag) Menus() []Menu {
    return append(t.Menus1, t.Menus2...)
}