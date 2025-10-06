package model

import (
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	menuID   uuid.UUIDs               `json:"menu_id" validate:"required,uuid4"`
	Quantity int                  `json:"quantity" validate:"required,min=1"`
	ApppointmentTime string        `json:"appointment_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	DeliveryMethod  string               `json:"delivery_method" validate:"required,oneof=pickup delivery"`
}
