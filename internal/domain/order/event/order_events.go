package event

import (
	"main/internal/domain/order/entity"
	"time"
)

// OrderEvent 订单事件基础结构
type OrderEvent struct {
	EventId   string
	OrderId   string
	Timestamp int64
	Type      string
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	OrderEvent
	Order *entity.Order
}

// NewOrderCreatedEvent 创建订单创建事件
func NewOrderCreatedEvent(order *entity.Order) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		OrderEvent: OrderEvent{
			EventId:   "", // 将由事件总线生成
			OrderId:   order.Id,
			Timestamp: time.Now().UnixMilli(),
			Type:      "OrderCreated",
		},
		Order: order,
	}
}

// OrderPaidEvent 订单支付事件
type OrderPaidEvent struct {
	OrderEvent
	Order       *entity.Order
	PaymentInfo interface{} // 支付信息
}

// NewOrderPaidEvent 创建订单支付事件
func NewOrderPaidEvent(order *entity.Order) *OrderPaidEvent {
	return &OrderPaidEvent{
		OrderEvent: OrderEvent{
			EventId:   "", // 将由事件总线生成
			OrderId:   order.Id,
			Timestamp: time.Now().UnixMilli(),
			Type:      "OrderPaid",
		},
		Order:       order,
		PaymentInfo: order.GetPaymentInfo(),
	}
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	OrderEvent
	Order *entity.Order
}

// NewOrderCancelledEvent 创建订单取消事件
func NewOrderCancelledEvent(order *entity.Order) *OrderCancelledEvent {
	return &OrderCancelledEvent{
		OrderEvent: OrderEvent{
			EventId:   "", // 将由事件总线生成
			OrderId:   order.Id,
			Timestamp: time.Now().UnixMilli(),
			Type:      "OrderCancelled",
		},
		Order: order,
	}
}

// OrderUpdatedEvent 订单更新事件
type OrderUpdatedEvent struct {
	OrderEvent
	Order *entity.Order
}

// NewOrderUpdatedEvent 创建订单更新事件
func NewOrderUpdatedEvent(order *entity.Order) *OrderUpdatedEvent {
	return &OrderUpdatedEvent{
		OrderEvent: OrderEvent{
			EventId:   "", // 将由事件总线生成
			OrderId:   order.Id,
			Timestamp: time.Now().UnixMilli(),
			Type:      "OrderUpdated",
		},
		Order: order,
	}
}
