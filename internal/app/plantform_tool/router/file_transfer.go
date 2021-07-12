package router

import (
	"github.com/gin-gonic/gin"
	"ls/internal/app/plantform_tool/handle"
)

func RegisteFileTransferRouter(engine *gin.Engine) {
	/*一、上传*/
	engine.POST("Upload", handle.File{}.Upload)
	/*二、下载*/
	engine.POST("Download", handle.File{}.Download)
	/*四、查询本地文件*/
	engine.POST("CheckExists", handle.File{}.CheckExists)
	/*十、读取文件*/
	engine.POST("ReadFromFile", handle.File{}.ReadFromFile)
	/*十一、压缩*/
	engine.POST("Compress", handle.File{}.Compress)
	/*十二、解压*/
	engine.POST("Decompress", handle.File{}.Decompress)

}