package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// ResponseHandler 统一响应处理中间件
func ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果已经有错误处理，则跳过
	if r.Response.BufferLength() > 0 {
		return
	}

	// 处理成功响应
	if data := r.GetHandlerResponse(); data != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    0,
			Message: "success",
			Data:    data,
		})
	}
}
