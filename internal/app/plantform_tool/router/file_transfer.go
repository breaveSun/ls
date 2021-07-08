package router

import (
	"github.com/gin-gonic/gin"
	"ls/internal/app/plantform_tool/handle"
)

func RegisteFileTransferRouter(engine *gin.Engine) {
	//扫描软件控制
	FileTransfer := engine.Group("fileTransfer")
	{
		//文件上传
		FileTransfer.POST("upload", handle.File{}.Upload)

		//文件下载
		FileTransfer.POST("download", handle.File{}.Download)

	}
}