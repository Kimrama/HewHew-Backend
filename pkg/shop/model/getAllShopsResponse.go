package model

type GetAllShopResponse struct {
	ShopID      string   `json:"shop_id"`
	Name        string   `json:"name"`
	CanteenName string   `json:"canteen_name"`
	Address     string   `json:"address"`
	ImageURL    string   `json:"shopimage_url"`
	State       bool     `json:"state"`
	Tags        []string `json:"tags"`
}
