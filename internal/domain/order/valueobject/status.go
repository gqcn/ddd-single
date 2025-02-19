package valueobject

// OrderStatus represents the status of an order
type OrderStatus string

const (
	// OrderStatusCreated represents a newly created order
	OrderStatusCreated OrderStatus = "created"
	// OrderStatusPaid represents a paid order
	OrderStatusPaid OrderStatus = "paid"
	// OrderStatusShipping represents an order that is being shipped
	OrderStatusShipping OrderStatus = "shipping"
	// OrderStatusDelivered represents a delivered order
	OrderStatusDelivered OrderStatus = "delivered"
	// OrderStatusCancelled represents a cancelled order
	OrderStatusCancelled OrderStatus = "cancelled"
)

// IsValid checks if the order status is valid
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusCreated, OrderStatusPaid, OrderStatusShipping,
		OrderStatusDelivered, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if the current status can transition to the target status
func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
	switch s {
	case OrderStatusCreated:
		return target == OrderStatusPaid || target == OrderStatusCancelled
	case OrderStatusPaid:
		return target == OrderStatusShipping || target == OrderStatusCancelled
	case OrderStatusShipping:
		return target == OrderStatusDelivered || target == OrderStatusCancelled
	case OrderStatusDelivered, OrderStatusCancelled:
		return false
	default:
		return false
	}
}

// String returns the string representation of the order status
func (s OrderStatus) String() string {
	return string(s)
}
