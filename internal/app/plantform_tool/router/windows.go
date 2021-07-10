package router

import (
	"github.com/gin-gonic/gin"
	"ls/internal/app/plantform_tool/handle"
)

func RegisteWindowsServerRouter(engine *gin.Engine) {
	engine.POST("ReadRegistry", handle.WindowsServer{}.ReadRegistry  )
}