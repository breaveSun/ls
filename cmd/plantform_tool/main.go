package main

import (
	"fmt"
	"go.uber.org/zap"
	"ls/internal/app/plantform_tool"
	"ls/internal/app/plantform_tool/lib"
	"ls/internal/app/plantform_tool/router"
	"ls/internal/pkg/lib/logger"
	"ls/internal/pkg/lib/redis"
	"net/http"
)
func main(){
	//初始化配置
	lib.InitConf()
	//初始化日志
	logger.InitLog()
	//错误捕获
	defer func() {
		err := recover()
		if err != nil {
			var restartCount =  plantform_tool.ServerConfig.Restart
			//0 不重启 -1一直重启  >0重启并且次数-1
			if restartCount == -1{
				//开启重启
				main()
			} else if restartCount > 0 {
				//重启次数-1
				plantform_tool.Config.Set("server_restart",restartCount-1)
				main()
			}
			logger.Logger.Error("服务异常", zap.String("err_info", fmt.Sprintf("%v", err)))
		}
	}()
	//初始化redis服务
	redis.PoolInitRedis(plantform_tool.RedisConfig)

	//初始化路由
	engine := router.InitRouter()

	//服务初始化
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", plantform_tool.ServerConfig.Port),
		Handler: engine,
	}
	/*//信号量检测 只能检测到停止
	go	app.SignGrab()*/
	//todo:端口检测
	//端口监听
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Fatal("监听端口失败，检查是否被占用", zap.Int("port", plantform_tool.ServerConfig.Port), zap.Reflect("error", err.Error()))
	}
	logger.Logger.Info("服务监听接口成功")


}