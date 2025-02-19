package event

import (
	"main/internal/domain/order/entity"
	"main/internal/infrastructure/eventbus"
)

const (
	OrderCreatedEventName       = "order.created"
	OrderCanceledEventName      = "order.canceled"
	OrderCompletedEventName     = "order.completed"
	OrderItemAddedEventName     = "order.item.added"
	OrderStatusChangedEventName = "order.status.changed"
)

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	eventbus.BaseEvent
	Order *entity.Order `json:"order"`
}

func NewOrderCreatedEvent(order *entity.Order) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		BaseEvent: eventbus.NewBaseEvent(OrderCreatedEventName, order.Id),
		Order:     order,
	}
}

// OrderCanceledEvent 订单取消事件
type OrderCanceledEvent struct {
	eventbus.BaseEvent
	OrderId string `json:"orderId"`
	Reason  string `json:"reason"`
}

func NewOrderCanceledEvent(orderId, reason string) *OrderCanceledEvent {
	return &OrderCanceledEvent{
		BaseEvent: eventbus.NewBaseEvent(OrderCanceledEventName, orderId),
		OrderId:   orderId,
		Reason:    reason,
	}
}

// OrderCompletedEvent 订单完成事件
type OrderCompletedEvent struct {
	eventbus.BaseEvent
	OrderId string `json:"orderId"`
}

func NewOrderCompletedEvent(orderId string) *OrderCompletedEvent {
	return &OrderCompletedEvent{
		BaseEvent: eventbus.NewBaseEvent(OrderCompletedEventName, orderId),
		OrderId:   orderId,
	}
}

// OrderItemAddedEvent 订单项添加事件
type OrderItemAddedEvent struct {
	eventbus.BaseEvent
	OrderId string            `json:"orderId"`
	Item    *entity.OrderItem `json:"item"`
}

func NewOrderItemAddedEvent(orderId string, item *entity.OrderItem) *OrderItemAddedEvent {
	return &OrderItemAddedEvent{
		BaseEvent: eventbus.NewBaseEvent(OrderItemAddedEventName, orderId),
		OrderId:   orderId,
		Item:      item,
	}
}

// OrderStatusChangedEvent 订单状态变更事件
type OrderStatusChangedEvent struct {
	eventbus.BaseEvent
	OrderId   string `json:"orderId"`
	OldStatus string `json:"oldStatus"`
	NewStatus string `json:"newStatus"`
}

func NewOrderStatusChangedEvent(orderId string, oldStatus, newStatus string) *OrderStatusChangedEvent {
	return &OrderStatusChangedEvent{
		BaseEvent: eventbus.NewBaseEvent(OrderStatusChangedEventName, orderId),
		OrderId:   orderId,
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}
}
