package model

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryMethod string

const (
	DeliveryMethodHandToHand DeliveryMethod = "handtohand"
	DeliveryMethodDropOff    DeliveryMethod = "dropoff"
)

type MenuItem struct {
	MenuID   uuid.UUID `json:"menu_id"`
	Quantity int       `json:"quantity"`
}

type CreateOrderRequest struct {
	Menu            []MenuItem     `json:"menu"`
	DropOffID       uuid.UUID      `json:"dropoff_location"`
	AppointmentTime time.Time      `json:"appointment_time"`
	DeliveryMethod  DeliveryMethod `json:"delivery_method"`
}
