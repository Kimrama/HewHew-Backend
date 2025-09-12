package model

type EditShopRequest struct {
    ShopName        string  `json:"shop_name"`
    ShopCanteenName string  `json:"canteen_name"`
}