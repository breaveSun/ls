package router

import (
	"github.com/gin-gonic/gin"
	"ls/internal/app/plantform_tool/handle"
)

func RegisteWindowsServerRouter(engine *gin.Engine) {
	/*五、读取注册表*/
	engine.POST("ReadRegistry", handle.WindowsServer{}.ReadRegistry)
	/*七-1、App 运行状态检测（一次性）*/
	engine.POST("CheckRunning", handle.WindowsServer{}.CheckRunning)
	/*七-2、App 运行状态检测（持续检测）*/
	engine.POST("RunningStatus", handle.WindowsServer{}.RunningStatus)
}