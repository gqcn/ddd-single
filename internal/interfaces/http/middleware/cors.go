package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// Cors 跨域中间件
func Cors(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
