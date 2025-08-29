package model

import (
	"hewhew-backend/utils"
)

type UserModel struct {
	Username string
	Password string
	FName    string
	LName    string
	Gender   string
	Image    *utils.ImageModel
}
