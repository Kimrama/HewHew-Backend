package model

type UserDetailResponse struct {
	Username        string `json:"username"`
	FName           string `json:"fname"`
	LName           string `json:"lname"`
	Gender          string `json:"gender"`
	ProfileImageURL string `json:"profile_image_url"`
}
