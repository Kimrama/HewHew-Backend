package model

import (
	"hewhew-backend/utils"
)

type CreateUserRequest struct {
	Username string
	Password string
	FName    string
	LName    string
	Gender   string
	Image    *utils.ImageModel
}
