package order

import (
	"time"

	"main/internal/domain/order/entity"
	"main/internal/domain/order/valueobject"

	"github.com/gogf/gf/v2/errors/gerror"
)

// OrderDTO represents the data transfer object for an order
type OrderDTO struct {
	Id          string          `json:"id"`
	UserId      string          `json:"userId"`
	Items       []*OrderItemDTO `json:"items"`
	TotalAmount float64         `json:"totalAmount"`
	Status      string          `json:"status"`
	Remark      string          `json:"remark,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// OrderItemDTO represents the data transfer object for an order item
type OrderItemDTO struct {
	Id          string  `json:"id"`
	ProductId   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	SubTotal    float64 `json:"subTotal"`
}

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	UserId string              `json:"userId" v:"required"`
	Items  []*OrderItemRequest `json:"items" v:"required|length:1,"`
	Remark string              `json:"remark,omitempty"`
}

// OrderItemRequest represents the request to create a new order item
type OrderItemRequest struct {
	ProductId   string  `json:"productId" v:"required"`
	ProductName string  `json:"productName" v:"required"`
	Quantity    int     `json:"quantity" v:"required|min:1"`
	Price       float64 `json:"price" v:"required|min:0.01"`
}

// UpdateOrderStatusRequest represents the request to update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" v:"required"`
	Remark string `json:"remark,omitempty"`
}

// ToDTO converts an Order entity to OrderDTO
func ToDTO(order *entity.Order) *OrderDTO {
	items := make([]*OrderItemDTO, len(order.Items))
	for i, item := range order.Items {
		items[i] = &OrderItemDTO{
			Id:          item.Id,
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price.Amount(),
			SubTotal:    item.SubTotal().Amount(),
		}
	}

	return &OrderDTO{
		Id:          order.Id,
		UserId:      order.UserId,
		Items:       items,
		TotalAmount: order.TotalAmount.Amount(),
		Status:      string(order.Status),
		Remark:      order.Remark,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

// ToDTOs converts a slice of Order entities to OrderDTOs
func ToDTOs(orders []*entity.Order) []*OrderDTO {
	dtos := make([]*OrderDTO, len(orders))
	for i, order := range orders {
		dtos[i] = ToDTO(order)
	}
	return dtos
}

// ToEntity converts CreateOrderRequest to Order entity
func (r *CreateOrderRequest) ToEntity() *entity.Order {
	order := entity.NewOrder(r.UserId)
	order.Remark = r.Remark

	for _, itemReq := range r.Items {
		item := entity.NewOrderItem(
			itemReq.ProductId,
			itemReq.ProductName,
			itemReq.Quantity,
			itemReq.Price,
		)
		order.AddItem(item)
	}

	return order
}

// ToOrderStatus converts string to OrderStatus
func ToOrderStatus(status string) (valueobject.OrderStatus, error) {
	orderStatus := valueobject.OrderStatus(status)
	if !orderStatus.IsValid() {
		return "", gerror.Newf("invalid order status: %s", status)
	}
	return orderStatus, nil
}
