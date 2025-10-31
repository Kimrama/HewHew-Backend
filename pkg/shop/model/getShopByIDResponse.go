package model

type MenuResponse struct {
	MenuID   string  `json:"menu_id"`
	Name     string  `json:"name"`
	Detail   string  `json:"detail"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
	ImageURL string  `json:"image_url"`
}

type GetShopByIdResponse struct {
	ShopID      string         `json:"shop_id"`
	Name        string         `json:"name"`
	CanteenName string         `json:"canteen_name"`
	Address     string         `json:"address"`
	ImageURL    string         `json:"shop_image_url"`
	State       bool           `json:"state"`
	Tags        []string       `json:"tags"`
	Menus       []MenuResponse `json:"menus"`
}
