package router

import (
	"github.com/gin-gonic/gin"
	"ls/internal/app/plantform_tool/handle"
)

func RegisteFileTransferRouter(engine *gin.Engine) {
/*	FileTransfer := engine.Group("fileTransfer")
	{*/
/*		//文件上传
	engine.POST("upload", handle.File{}.Upload)

		//文件下载
	engine.POST("download", handle.File{}.Download)*/
	/*四、查询本地文件*/
	engine.POST("CheckExists", handle.File{}.CheckExists)
	/*十、读取文件*/
	engine.POST("ReadFromFile", handle.File{}.ReadFromFile)
	/*十一、压缩*/
	engine.POST("Compress", handle.File{}.Compress)
	/*十二、解压*/
	engine.POST("Decompress", handle.File{}.Decompress)

}