package router

import (
	"github.com/gin-gonic/gin"
	"ls/internal/app/plantform_tool/handle"
)

func RegisterControlSoftwareRouter(engine *gin.Engine){
	//获取软件版本信息
	//engine.POST("getSoftwareInfo","")
	//扫描软件控制
	scanControl := engine.Group("scan")
	{
		//开始扫描
		scanControl.POST("start",handle.Scan{}.ScanStart)
		//获取scan参数
		//完成scan
		scanControl.POST("over",handle.Scan{}.ScanOver)
		//下载scan
		//查看扫描状态
	}
	//口内扫描软件控制路由组

	//CAD软件控制路由组

	//CAM软件控制路由组

	//CNC软件控制路由组

}