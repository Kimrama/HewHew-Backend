package model

type CreateNotificationRequest struct {
	ReceiverID string `json:"receiver_id"`
	OrderID    string `json:"order_id"`
	Topic      string `json:"topic"`
	Message    string `json:"message"`
}