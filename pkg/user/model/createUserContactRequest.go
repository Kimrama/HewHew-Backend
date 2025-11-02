package model

type EditUserContactRequest struct {
	ContactType string `json:"contact_type"`
	Detail      string `json:"contact_detail"`
}
