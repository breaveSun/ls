package router

import (
	"github.com/gin-gonic/gin"
	"ls/internal/app/plantform_tool/handle"
)

func RegisteCommandRouter(engine *gin.Engine) {
	/*六、CMD 执行*/
	engine.POST("ExecCommand", handle.Command{}.ExecCommand)
	/*八、杀进程*/
	engine.POST("KillProcess", handle.Command{}.KillProcess)
}