package model

type EditMenuRequest struct {
	Name   *string `json:"name,omitempty"`
	Detail *string `json:"detail,omitempty"`
	Price  *string `json:"price,omitempty"`
	Tag1ID *string `json:"tag1_id,omitempty"`
	Tag2ID *string `json:"tag2_id,omitempty"`
}
