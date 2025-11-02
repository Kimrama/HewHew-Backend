package model

import "github.com/google/uuid"

type EditUserRequest struct {
	UserID uuid.UUID `json:"user_id"`
	FName  *string   `json:"fname,omitempty"`
	LName  *string   `json:"lname,omitempty"`
	Gender *string   `json:"gender,omitempty"`
}
