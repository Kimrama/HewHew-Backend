package model

import (
	"github.com/google/uuid"
)

type AcceptOrderRequest struct {
	DeliveryuserID uuid.UUID
	OrderID        uuid.UUID
}
