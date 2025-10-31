package model

import (
	"github.com/google/uuid"
)

type GetReviewResponse struct {
	UserTargetID uuid.UUID `json:"user_target_id"`
	OrderID      uuid.UUID `json:"order_id"`
	Rating       int       `json:"rating"`
	Comment      string    `json:"comment"`
}
