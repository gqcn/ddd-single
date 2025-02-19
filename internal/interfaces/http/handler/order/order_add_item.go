package order

import (
	"context"

	"main/internal/application/order/dto"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// AddItemReq 添加订单项请求
type AddItemReq struct {
	g.Meta `path:"/orders/{id}/items" method:"post" tags:"订单" summary:"添加订单项"`
	Id     string `v:"required" path:"id" dc:"订单Id"`
	dto.OrderItemRequest
}

// AddItemRes 添加订单项响应
type AddItemRes struct{}

// AddItem 添加订单项
func (o *Order) AddItem(ctx context.Context, req *AddItemReq) (res *AddItemRes, err error) {
	if err := o.orderApp.AddOrderItem(ctx, req.Id, &req.OrderItemRequest); err != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, err.Error())
	}
	return &AddItemRes{}, nil
}
