package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateReviewRequest struct {
	Menu            []MenuItem     `json:"menu"`
	DropOffID       uuid.UUID      `json:"dropoff_location"`
	AppointmentTime time.Time      `json:"appointment_time"`
	DeliveryMethod  DeliveryMethod `json:"delivery_method"`
}
