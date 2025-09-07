package model

type CanteenRequest struct {
	CanteenName string `json:"canteen_name"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}