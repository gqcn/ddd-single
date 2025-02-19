package order

import (
	"context"

	"main/internal/application/order/dto"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// UpdateStatusReq 更新订单状态请求
type UpdateStatusReq struct {
	g.Meta `path:"/orders/{id}/status" method:"put" tags:"订单" summary:"更新订单状态"`
	Id     string `v:"required" path:"id" dc:"订单Id"`
	dto.UpdateOrderStatusRequest
}

// UpdateStatusRes 更新订单状态响应
type UpdateStatusRes struct{}

// UpdateStatus 更新订单状态
func (o *Order) UpdateStatus(ctx context.Context, req *UpdateStatusReq) (res *UpdateStatusRes, err error) {
	if err := o.orderApp.UpdateOrderStatus(ctx, req.Id, &req.UpdateOrderStatusRequest); err != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, err.Error())
	}
	return &UpdateStatusRes{}, nil
}
