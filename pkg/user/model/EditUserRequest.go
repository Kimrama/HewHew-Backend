package model

import "github.com/google/uuid"

type EditUserRequest struct {
	UserID uuid.UUID `json:"user_id"`
	FName  string    `json:"fname"`
	LName  string    `json:"lname"`
	Gender string    `json:"gender"`
}
