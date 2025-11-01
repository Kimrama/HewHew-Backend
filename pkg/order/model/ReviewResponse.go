package model

import (
	"time"

	"github.com/google/uuid"
)

type GetReviewResponse struct {
	ReviewID       uuid.UUID `json:"review_id"`
	UserReviewerID uuid.UUID `json:"user_reviewer_id"`
	UserTargetID   uuid.UUID `json:"user_target_id"`
	OrderID        uuid.UUID `json:"order_id"`
	Rating         int       `json:"rating"`
	Comment        string    `json:"comment"`
	TimeStamp      time.Time `json:"time_stamp"`
}
