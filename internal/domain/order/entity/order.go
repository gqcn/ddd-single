package entity

import (
	"time"

	"main/internal/domain/order/valueobject"
	sharedvo "main/internal/domain/shared/valueobject"

	"github.com/gogf/gf/v2/errors/gerror"
)

// Order represents the order aggregate root
type Order struct {
	Id           string
	UserId       string
	Items        []*OrderItem
	TotalAmount  *sharedvo.Money
	Status       valueobject.OrderStatus
	PaymentInfo  *valueobject.PaymentInfo // 支付信息
	Remark       string
	CreatedAt    int64
	UpdatedAt    int64
	PaidAt       int64 // 支付时间
}

// NewOrder creates a new order instance
func NewOrder(userId string) *Order {
	return &Order{
		Id:          "", // ID will be assigned by the infrastructure layer
		UserId:      userId,
		Status:      valueobject.OrderStatusCreated,
		Items:       make([]*OrderItem, 0),
		TotalAmount: sharedvo.NewMoney(0, "CNY"),
		CreatedAt:   time.Now().UnixMilli(),
		UpdatedAt:   time.Now().UnixMilli(),
		PaidAt:      0,
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

// ProcessPayment 处理订单支付
// 这是一个领域行为，包含了支付相关的业务规则
func (o *Order) ProcessPayment(paymentInfo *valueobject.PaymentInfo) error {
	// 1. 验证订单状态
	if o.Status != valueobject.OrderStatusCreated {
		return gerror.Newf("cannot pay order in status: %s", o.Status)
	}

	// 2. 验证支付金额
	if !o.TotalAmount.Equals(paymentInfo.Amount) {
		return gerror.New("payment amount does not match order amount")
	}

	// 3. 更新订单状态和支付信息
	if err := o.UpdateStatus(valueobject.OrderStatusPaid); err != nil {
		return gerror.Wrap(err, "failed to update order status")
	}

	o.PaymentInfo = paymentInfo
	o.PaidAt = time.Now().UnixMilli()
	o.UpdatedAt = o.PaidAt

	return nil
}

// GetPaymentInfo 获取支付信息
func (o *Order) GetPaymentInfo() *valueobject.PaymentInfo {
	return o.PaymentInfo
}

// IsPaid 检查订单是否已支付
func (o *Order) IsPaid() bool {
	return o.Status == valueobject.OrderStatusPaid && o.PaymentInfo != nil
}

// Cancel 取消订单
// 这是一个领域行为，包含了取消订单的业务规则
func (o *Order) Cancel() error {
	// 1. 验证订单是否可以取消
	if o.Status != valueobject.OrderStatusCreated {
		return gerror.Newf("cannot cancel order in status: %s", o.Status)
	}

	// 2. 更新订单状态
	if err := o.UpdateStatus(valueobject.OrderStatusCancelled); err != nil {
		return gerror.Wrap(err, "failed to update order status")
	}

	// 3. 更新时间戳
	o.UpdatedAt = time.Now().UnixMilli()

	return nil
}

// IsCancelled 检查订单是否已取消
func (o *Order) IsCancelled() bool {
	return o.Status == valueobject.OrderStatusCancelled
}

// recalculateTotal recalculates the total amount of the order
func (o *Order) recalculateTotal() {
	total := sharedvo.NewMoney(0, "CNY")
	for _, item := range o.Items {
		itemTotal := item.Price.Multiply(float64(item.Quantity))
		newTotal, _ := total.Add(itemTotal)
		total = newTotal
	}
	o.TotalAmount = total
}

// UpdateRemark 更新订单备注
func (o *Order) UpdateRemark(remark string) {
	o.Remark = remark
	o.UpdatedAt = time.Now().UnixMilli()
}

// Validate 验证订单
func (o *Order) Validate() error {
	// 1. 基本字段验证
	if o.UserId == "" {
		return gerror.New("user id is required")
	}

	if len(o.Items) == 0 {
		return gerror.New("order must have at least one item")
	}

	if !o.Status.IsValid() {
		return gerror.Newf("invalid order status: %s", o.Status)
	}

	// 2. 验证订单项
	for _, item := range o.Items {
		if err := item.Validate(); err != nil {
			return gerror.Wrap(err, "invalid order item")
		}
	}

	// 3. 如果是已支付状态，验证支付信息
	if o.Status == valueobject.OrderStatusPaid {
		if o.PaymentInfo == nil {
			return gerror.New("payment info is required for paid order")
		}
		if err := o.PaymentInfo.Validate(); err != nil {
			return gerror.Wrap(err, "invalid payment info")
		}
	}

	// 4. 验证时间戳
	if o.CreatedAt <= 0 {
		return gerror.New("invalid created time")
	}
	if o.UpdatedAt < o.CreatedAt {
		return gerror.New("updated time cannot be earlier than created time")
	}

	return nil
}
