package model

import (
	"hewhew-backend/utils"
)

type MenuStatus string

const (
	MenuStatusAvailable   MenuStatus = "available"
	MenuStatusUnavailable MenuStatus = "unavailable"
)

type CreateMenuRequest struct {
	Name   string           `json:"name"`
	Detail string           `json:"detail"`
	Price  string           `json:"price"`
	Status MenuStatus       `json:"status"`
	Tag1ID string            `json:"tag1_id"`
	Tag2ID string            `json:"tag2_id"`
	Image  *utils.ImageModel `json:"image"`
}

