package order

import (
	"main/internal/application/order"
)

// Order 订单控制器
type Order struct {
	orderApp *order.ApplicationService
}

// NewOrder 创建订单控制器实例
func NewOrder(orderApp *order.ApplicationService) *Order {
	return &Order{
		orderApp: orderApp,
	}
}
