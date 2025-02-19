package order

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// CancelReq 取消订单请求
type CancelReq struct {
	g.Meta `path:"/orders/{id}/cancel" method:"post" tags:"订单" summary:"取消订单"`
	Id     string `v:"required" path:"id" dc:"订单Id"`
	Reason string `v:"required" json:"reason" dc:"取消原因"`
}

// CancelRes 取消订单响应
type CancelRes struct{}

// Cancel 取消订单
func (o *Order) Cancel(ctx context.Context, req *CancelReq) (res *CancelRes, err error) {
	if err := o.orderApp.CancelOrder(ctx, req.Id, req.Reason); err != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, err.Error())
	}
	return &CancelRes{}, nil
}
