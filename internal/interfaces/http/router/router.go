package router

import (
	"main/internal/interfaces/api/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Register 注册所有路由
func Register(server *ghttp.Server) {
	// 注册全局中间件
	server.Use(
		middleware.ErrorHandler,
		middleware.ResponseHandler,
		middleware.Cors,
	)

	// 注册 API 路由组
	server.Group("/api/v1", func(group *ghttp.RouterGroup) {
		// 添加认证中间件
		group.Middleware(middleware.Auth)

		// 注册模块路由
		registerOrderRoutes(group)
		// TODO: 注册其他模块路由
	})

	// 注册 OpenAPI 路由
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/api.json", func(r *ghttp.Request) {
			r.Response.WriteJson(g.Map{
				"openapi": "3.0.0",
				"info": g.Map{
					"title":   "电商系统 API",
					"version": "v1",
				},
				// TODO: 完善 OpenAPI 文档
			})
		})
	})
}
