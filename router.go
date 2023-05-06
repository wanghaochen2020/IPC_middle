package main

import (
	"github.com/gin-gonic/gin"
	"middle/middleware"
	"middle/routers"
	"net/http"
)

func InitRouter() { //可以返回一个*gin.Engine
	gin.SetMode("debug")
	r := gin.New()
	// 加载中间件
	r.Use(middleware.Cors())     // 跨域中间件
	r.Use(gin.Recovery())        // 恢复恐慌中间件
	_ = r.SetTrustedProxies(nil) // 信任所有ip端口

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该路由",
		})
	})
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该方法",
		})
	})

	// 获取气象站数据路由
	routers.DeviceRouter(r)
	// 获取垃圾桶数据路由
	//routers.TrashCanRouter(r)

	_ = r.Run(":7766")
}
