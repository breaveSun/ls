package handle

import (
	"github.com/antage/eventsource"
	"github.com/gin-gonic/gin"
	form2 "ls/internal/app/plantform_tool/form"
	"ls/internal/pkg/common"
	"time"
)

type File struct {
	common.BaseHandler
}
//上传
func (h *File) Upload(c *gin.Context){
	//接收参数
	var form form2.UploadFileForm
	if err := h.BindParams(c, &form); err != nil {
		h.HandleError(c, err)
		return
	}
	// 判断本地是否存在该文件
	// 存在 ：判断版本是否为最新
	// 不存在｜不是最新版本则
	// 并保存到redis
	//创建一个sse链接
	es := eventsource.New(
		&eventsource.Settings{
			Timeout: 5 * time.Second,
			CloseOnTimeout: false,
			//空闲链接超时时间
			IdleTimeout: 30 * time.Minute,
		}, nil)
	es.SendRetryMessage(3 * time.Second)
	//defer es.Close()

}
//下载
func (h *File) Download(c *gin.Context){

}