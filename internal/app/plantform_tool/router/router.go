package router

import (
	"github.com/gin-gonic/gin"
	"ls/internal/pkg/lib/logger"
)

func InitRouter() *gin.Engine {
	engine := gin.Default()
	RegisteFileTransferRouter(engine)
	RegisteWindowsServerRouter(engine)
	RegisteCommandRouter(engine)

	logger.Logger.Info("路由初始话成功")
	return engine
}
