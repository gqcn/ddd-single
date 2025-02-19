package order

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"main/internal/application/order/dto"
)

// CreateReq 创建订单请求
type CreateReq struct {
	g.Meta `path:"/orders" method:"post" tags:"订单" summary:"创建订单"`
	dto.CreateOrderRequest
}

// CreateRes 创建订单响应
type CreateRes struct {
	*dto.OrderDTO
}

// Create 创建订单
func (o *Order) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	order, err := o.orderApp.CreateOrder(ctx, &req.CreateOrderRequest)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, err.Error())
	}
	return &CreateRes{OrderDTO: order}, nil
}
