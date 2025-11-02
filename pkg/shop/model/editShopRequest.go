package model

type EditShopRequest struct {
	ShopName        *string `json:"shop_name,omitempty"`
	ShopCanteenName *string `json:"canteen_name,omitempty"`
}
