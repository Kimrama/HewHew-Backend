package model

import "github.com/google/uuid"

type GetMenuByIDResponse struct {
	MenuID   uuid.UUID  `json:"menu_id"`
	Name     string     `json:"name"`
	Detail   string     `json:"detail"`
	Price    float64    `json:"price"`
	Status   string     `json:"status"`
	ImageURL string     `json:"image_url"`
	Tag1ID   *uuid.UUID `json:"tag1_id,omitempty"`
	Tag2ID   *uuid.UUID `json:"tag2_id,omitempty"`
}
