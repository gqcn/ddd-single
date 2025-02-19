package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"main/internal/interfaces/api/router"
)

func main() {
	var (
		ctx = gctx.New()
		app = g.Server()
	)

	// 注册路由
	router.Register(app)

	// 启动服务
	app.Run()
