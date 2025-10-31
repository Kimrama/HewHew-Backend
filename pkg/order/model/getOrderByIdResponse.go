package model

import (
	"time"

	"github.com/google/uuid"
)

type MenuQuantityResponse struct {
	MenuID   uuid.UUID `json:"menu_id"`
	Quantity int       `json:"quantity"`
}

type GetOrderByIdResponse struct {
	OrderID              uuid.UUID              `json:"order_id"`
	UserOrderID          uuid.UUID              `json:"user_order_id"`
	UserDeliveryID       *uuid.UUID             `json:"user_delivery_id,omitempty"`
	Status               string                 `json:"status"`
	OrderDate            time.Time              `json:"order_date"`
	DeliveryMethod       string                 `json:"delivery_method"`
	ConfirmationImageURL string                 `json:"confirmation_image_url,omitempty"`
	AppointmentTime      time.Time              `json:"appointment_time"`
	DropOffLocationID    uuid.UUID              `json:"drop_off_location_id"`
	MenuQuantity         []MenuQuantityResponse `json:"menu_quantity"`
	ShopName             string                 `json:"shop_name"`
	CanteenName          string                 `json:"canteen_name"`
	ShippingFee          float64                `json:"shipping_fee"`
}
