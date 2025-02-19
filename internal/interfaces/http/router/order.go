package router

import (
	"main/internal/application/order"
	"main/internal/domain/order/service"
	"main/internal/infrastructure/persistence/mysql"
	orderHandler "main/internal/interfaces/http/handler/order"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// registerOrderRoutes 注册订单相关路由
func registerOrderRoutes(group *ghttp.RouterGroup) {
	// 初始化依赖
	db := g.DB()
	orderRepo := mysql.NewOrderRepository(db)
	orderDomainService := service.NewOrderService(orderRepo, productService)
	orderApp := order.NewApplicationService(orderRepo, orderDomainService)

	// 创建处理器
	handler := orderHandler.NewOrder(orderApp)

	// 注册路由
	group.Group("/orders", func(group *ghttp.RouterGroup) {
		// 创建订单
		group.POST("/", handler.Create)

		// 获取订单详情
		group.GET("/{id}", handler.Get)

		// 更新订单状态
		group.PUT("/{id}/status", handler.UpdateStatus)

		// 取消订单
		group.POST("/{id}/cancel", handler.Cancel)

		// 添加订单项
		group.POST("/{id}/items", handler.AddItem)
	})

	// 用户订单列表
	group.GET("/users/{userId}/orders", handler.List)
}
