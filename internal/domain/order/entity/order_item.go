package entity

import (
	sharedvo "main/internal/domain/shared/valueobject"

	"github.com/gogf/gf/v2/errors/gerror"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ProductId   string          // 商品ID
	ProductName string          // 商品名称
	Quantity    int             // 数量
	Price       *sharedvo.Money // 单价
}

// NewOrderItem creates a new order item
func NewOrderItem(productId string, productName string, quantity int, price float64) *OrderItem {
	return &OrderItem{
		ProductId:   productId,
		ProductName: productName,
		Quantity:    quantity,
		Price:       sharedvo.NewMoney(price, "CNY"),
	}
}

// GetSubtotal calculates the subtotal for this item
func (i *OrderItem) GetSubtotal() *sharedvo.Money {
	return i.Price.Multiply(float64(i.Quantity))
}

// UpdateQuantity updates the quantity of the item
func (i *OrderItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return gerror.New("quantity must be positive")
	}
	i.Quantity = quantity
	return nil
}

// Validate validates the order item
func (i *OrderItem) Validate() error {
	if i.ProductId == "" {
		return gerror.New("product id is required")
	}

	if i.ProductName == "" {
		return gerror.New("product name is required")
	}

	if i.Quantity <= 0 {
		return gerror.New("quantity must be positive")
	}

	if i.Price == nil {
		return gerror.New("price is required")
	}

	if i.Price.Amount() <= 0 {
		return gerror.New("price must be positive")
	}

	return nil
}
