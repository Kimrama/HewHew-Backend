
package service

import (
    "hewhew-backend/pkg/order/repository"
)

type OrderServiceImpl struct {
    OrderRepository repository.OrderRepository
}

func NewOrderServiceImpl(OrderRepository repository.OrderRepository) OrderService {
    return &OrderServiceImpl{
        OrderRepository: OrderRepository,
    }
}
