package model

import "time"

type GetNotification struct {
	NotificationID string    `json:"notification_id"`
	OrderID        string    `json:"order_id"`
	ReceiverID     string    `json:"receiver_id"`
	Topic          string    `json:"topic"`
	Message        string    `json:"message"`
	TimeStamp      time.Time `json:"time_stamp"`
}