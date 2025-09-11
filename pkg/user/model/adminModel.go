package model

type CreateAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FName    string `json:"fname"`
	LName    string `json:"lname"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	CanteenName string `json:"canteen_name"`
}
