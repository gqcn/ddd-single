package entity

import (
	"time"

	"main/internal/domain/order/valueobject"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/guid"
)

// Order represents the order aggregate root
type Order struct {
	Id          string
	UserId      string
	Items       []*OrderItem
	TotalAmount *valueobject.Money
	Status      valueobject.OrderStatus
	Remark      string
	CreatedAt   int64
	UpdatedAt   int64
}

// NewOrder creates a new order instance
func NewOrder(userId string) *Order {
	return &Order{
		Id:          guid.S(),
		UserId:      userId,
		Status:      valueobject.OrderStatusCreated,
		Items:       make([]*OrderItem, 0),
		TotalAmount: valueobject.NewMoney(0, "CNY"),
		CreatedAt:   time.Now().UnixMilli(),
		UpdatedAt:   time.Now().UnixMilli(),
	}
}

// AddItem adds a new item to the order
func (o *Order) AddItem(item *OrderItem) error {
	if o.Status != valueobject.OrderStatusCreated {
		return gerror.New("cannot add items to non-created order")
	}

	// Check if product already exists in order
	for _, existingItem := range o.Items {
		if existingItem.ProductId == item.ProductId {
			existingItem.Quantity += item.Quantity
			o.recalculateTotal()
			o.UpdatedAt = time.Now().UnixMilli()
			return nil
		}
	}

	o.Items = append(o.Items, item)
	o.recalculateTotal()
	o.UpdatedAt = time.Now().UnixMilli()
	return nil
}

// RemoveItem removes an item from the order
func (o *Order) RemoveItem(productId string) error {
	if o.Status != valueobject.OrderStatusCreated {
		return gerror.New("cannot remove items from non-created order")
	}

	for i, item := range o.Items {
		if item.ProductId == productId {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.recalculateTotal()
			o.UpdatedAt = time.Now().UnixMilli()
			return nil
		}
	}

	return gerror.Newf("product %s not found in order", productId)
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(newStatus valueobject.OrderStatus) error {
	if !newStatus.IsValid() {
		return gerror.Newf("invalid status: %s", newStatus)
	}

	if !o.Status.CanTransitionTo(newStatus) {
		return gerror.Newf("invalid status transition from %s to %s", o.Status, newStatus)
	}

	o.Status = newStatus
	o.UpdatedAt = time.Now().UnixMilli()
	return nil
}

// recalculateTotal recalculates the total amount of the order
func (o *Order) recalculateTotal() {
	total := valueobject.NewMoney(0, "CNY")
	for _, item := range o.Items {
		itemTotal := item.Price.Multiply(float64(item.Quantity))
		newTotal, _ := total.Add(itemTotal)
		total = newTotal
	}
	o.TotalAmount = total
}

// Validate validates the order
func (o *Order) Validate() error {
	if o.UserId == "" {
		return gerror.New("user Id is required")
	}

	if len(o.Items) == 0 {
		return gerror.New("order must contain at least one item")
	}

	for _, item := range o.Items {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
