package model

type TagCreateRequest struct {
	TagID  string `json:"tag_id"`
	ShopID string `json:"shop_id"`
	Topic  string `json:"topic"`
}