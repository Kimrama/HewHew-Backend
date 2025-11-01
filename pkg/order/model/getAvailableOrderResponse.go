package model

import (
	"time"

	"github.com/google/uuid"
)

type GetAvailableOrderResponse struct {
	OrderID           uuid.UUID              `json:"order_id"`
	UserOrderID       uuid.UUID              `json:"user_order_id"`
	Status            string                 `json:"status"`
	OrderDate         time.Time              `json:"order_date"`
	DeliveryMethod    string                 `json:"delivery_method"`
	AppointmentTime   time.Time              `json:"appointment_time"`
	DropOffLocationID uuid.UUID              `json:"drop_off_location_id"`
	MenuQuantity      []MenuQuantityResponse `json:"menu_quantity"`
	Amount            float64                `json:"amount"`
	ShopName          string                 `json:"shop_name"`
	ShopID            uuid.UUID              `json:"shop_id"`
	CanteenName       string                 `json:"canteen_name"`
	ShippingFee       float64                `json:"shipping_fee"`
}
