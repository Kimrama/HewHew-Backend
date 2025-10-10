package model

import "hewhew-backend/utils"

type CreateDropOffRequest struct {
	Latitude  string
	Longitude string
	Name      string
	Detail    string
	Image     *utils.ImageModel
}
