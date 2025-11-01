package model

type GetAllTagsRequest struct {
	Topic  string `json:"topic"`
	TagID  string `json:"tag_id"`
}