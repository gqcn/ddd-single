package entity

import (
	"main/internal/domain/order/valueobject"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gutil"
)

// OrderItem represents an item in an order
type OrderItem struct {
	Id          string
	OrderId     string
	ProductId   string
	ProductName string
	Quantity    int
	Price       *valueobject.Money
}

// NewOrderItem creates a new order item
func NewOrderItem(productId string, productName string, quantity int, price float64) *OrderItem {
	return &OrderItem{
		Id:          gutil.UUId(),
		ProductId:   productId,
		ProductName: productName,
		Quantity:    quantity,
		Price:       valueobject.NewMoney(price, "CNY"),
	}
}

// UpdateQuantity updates the quantity of the order item
func (i *OrderItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return gerror.New("quantity must be greater than 0")
	}
	i.Quantity = quantity
	return nil
}

// SubTotal calculates the subtotal for this item
func (i *OrderItem) SubTotal() *valueobject.Money {
	return i.Price.Multiply(float64(i.Quantity))
}

// Validate validates the order item
func (i *OrderItem) Validate() error {
	if i.ProductId == "" {
		return gerror.New("product Id is required")
	}

	if i.ProductName == "" {
		return gerror.New("product name is required")
	}

	if i.Quantity <= 0 {
		return gerror.New("quantity must be greater than 0")
	}

	if i.Price == nil || i.Price.Amount() <= 0 {
		return gerror.New("price must be greater than 0")
	}

	return nil
}
