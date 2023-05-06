package main

import "middle/model"

func main() {
	model.InitMongo()
	// 引用路由组件
	InitRouter()
}
