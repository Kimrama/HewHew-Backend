package model

import (
	"hewhew-backend/utils"
)

type CreateUserRequest struct {
	Username      string            `json:"username"`
	Password      string            `json:"password"`
	FName         string            `json:"fname"`
	LName         string            `json:"lname"`
	Gender        string            `json:"gender"`
	Image         *utils.ImageModel `json:"image,omitempty"`
	ContactType   string            `json:"contact_type"`
	ContactDetail string            `json:"contact_detail"`
}
