package model

type DropOffDetailResponse struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Name      string `json:"name"`
	Detail    string `json:"detail"`
	Image     string `json:"image_url"`
}
