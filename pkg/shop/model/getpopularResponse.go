package model

type ShopPopularResponse struct {
	ShopID      string         `json:"shop_id"`
	Name        string         `json:"name"`
	CanteenName string         `json:"canteen_name"`
	Address     string         `json:"address"`
	ImageURL    string         `json:"shop_image_url"`
	State       bool           `json:"state"`
	Tags        []string       `json:"tags"`
	Menus       []MenuResponse `json:"menus"`
}