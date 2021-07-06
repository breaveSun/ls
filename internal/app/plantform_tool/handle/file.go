package handle

import (
	"fmt"
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
	var form form2.UploadFileForm
	if err := h.BindParams(c, &form); err != nil {
		fmt.Println(err)
		//h.HandleError(c, err)
		return
	}
	//接收参数和返回参数
	es := eventsource.New(
		&eventsource.Settings{
			Timeout: 5 * time.Second,
			CloseOnTimeout: false,
			//空闲链接超时时间
			IdleTimeout: 30 * time.Minute,
		}, nil)
	es.SendRetryMessage(3 * time.Second)
	defer es.Close()
}
//下载
func (h *File) Download(c *gin.Context){

}