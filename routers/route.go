package routers

import (
	"github.com/gin-gonic/gin"
	"middle/api"
)

// 中转路由
func DeviceRouter(router *gin.Engine) {
	router.GET("api/getData", api.GetData) //
}
