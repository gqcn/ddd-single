package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Auth 认证中间件
func Auth(r *ghttp.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "未授权访问"))
		return
	}

	// TODO: 实现token验证逻辑

	r.Middleware.Next()
}
