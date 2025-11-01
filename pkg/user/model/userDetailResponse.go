package model

import "github.com/google/uuid"

type UserDetailResponse struct {
	UserID          uuid.UUID `json:"user_id"`
	Username        string    `json:"username"`
	FName           string    `json:"fname"`
	LName           string    `json:"lname"`
	Gender          string    `json:"gender"`
	ProfileImageURL string    `json:"profile_image_url"`
}
