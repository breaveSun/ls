package main

import (
	"fmt"
	"ls/internal/app/plantform_tool"
	"ls/internal/app/plantform_tool/lib"
	"ls/internal/app/plantform_tool/router"
)
func main(){
	//初始化配置
	lib.InitConf()
	if plantform_tool.RestartConfig == 0{
		//开启重启
	}
	//端口检测
	//端口被占用kill之后重启
	//初始化日志
	lib.InitLog()
	//初始化redis服务
	lib.PoolInitRedis()
	//初始化路由
	engine := router.InitRouter()
	err := engine.Run(":22456")
	if err != nil{
		fmt.Println("server err")
	}

}