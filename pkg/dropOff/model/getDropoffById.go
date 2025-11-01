package model

type GetDropOffByIDResponse struct {
	DropOffID string `json:"dropoff_id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Name      string `json:"name"`
	Detail    string `json:"detail"`
	Image     string `json:"image_url"`
}
