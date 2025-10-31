package model

import (
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type ConfirmOrderRequest struct {
	OrderID uuid.UUID
	Image   *utils.ImageModel
}
