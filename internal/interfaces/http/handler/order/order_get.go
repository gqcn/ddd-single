package order

import (
	"context"

	"main/internal/application/order/dto"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// GetReq 获取订单请求
type GetReq struct {
	g.Meta `path:"/orders/{id}" method:"get" tags:"订单" summary:"获取订单详情"`
	Id     string `v:"required" path:"id" dc:"订单Id"`
}

// GetRes 获取订单响应
type GetRes struct {
	*dto.OrderDTO
}

// Get 获取订单详情
func (o *Order) Get(ctx context.Context, req *GetReq) (res *GetRes, err error) {
	order, err := o.orderApp.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, err.Error())
	}
	return &GetRes{OrderDTO: order}, nil
}
