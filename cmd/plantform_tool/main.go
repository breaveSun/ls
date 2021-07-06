package main

import (
	"fmt"
	"ls/internal/app/plantform_tool/lib"
	"ls/internal/app/plantform_tool/router"
)
func main(){
	//初始化日志
	lib.InitLog()
	//初始化redis服务
	lib.InitRedis()
	//初始化路由
	engine := router.InitRouter()
	err := engine.Run(":22456")
	if err != nil{
		fmt.Println("server err")
	}
	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	//r.Run() // listen and serve on 0.0.0.0:8080
}