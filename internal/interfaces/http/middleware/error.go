package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()

	if err := r.GetError(); err != nil {
		code := gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}

		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    code.Code(),
			Message: err.Error(),
		})
	}
}
