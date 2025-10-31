package model

type GetShopResponse struct {
	ShopName string `json:"shop_name"`
	ShopID   string `json:"shop_id"`
}

type GetCanteenResponse struct {
	CanteenName string            `json:"canteen_name"`
	Latitude    string            `json:"latitude"`
	Longitude   string            `json:"longitude"`
	Shop        []GetShopResponse `json:"shop"`
}
