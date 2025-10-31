package model

type GetAllShopRequest struct {
    ShopName        string  `json:"shop_name"`
    ShopCanteenName string  `json:"shop_canteen_name"`
    ShopImg         string  `json:"shop_img"`
    State           bool    `json:"state"`
}