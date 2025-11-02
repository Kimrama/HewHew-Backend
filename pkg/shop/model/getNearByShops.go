package model

type GetNearByShopResponse struct {
	ShopID      string `json:"shop_id"`
	Name        string `json:"name"`
	CanteenName string `json:"canteen_name"`
	ImageURL    string `json:"shop_image_url"`
}

type GetNearByShopRequest struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
