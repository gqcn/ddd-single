package order

import (
	"context"

	"main/internal/application/order/dto"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// ListReq 获取用户订单列表请求
type ListReq struct {
	g.Meta  `path:"/users/{userId}/orders" method:"get" tags:"订单" summary:"获取用户订单列表"`
	UserId  string `v:"required" path:"userId" dc:"用户Id"`
	Page    int    `v:"min:1" query:"page" dc:"页码" d:"1"`
	PerPage int    `v:"min:1,max:50" query:"perPage" dc:"每页数量" d:"20"`
}

// ListRes 获取用户订单列表响应
type ListRes struct {
	List  []*dto.OrderDTO `json:"list"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
}

// List 获取用户订单列表
func (o *Order) List(ctx context.Context, req *ListReq) (res *ListRes, err error) {
	orders, err := o.orderApp.GetUserOrders(ctx, req.UserId)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, err.Error())
	}
	return &ListRes{
		List:  orders,
		Total: len(orders),
		Page:  req.Page,
	}, nil
}
