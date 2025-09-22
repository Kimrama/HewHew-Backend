package model

type ShopByAdminID struct {
    ShopName        string  `json:"shopname"`
    ShopCanteenName string  `json:"shopcanteenname"`
    ShopImg         string  `json:"shopimg"`
    State           bool    `json:"state"`
}