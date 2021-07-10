package handle

import (
	"fmt"
	"github.com/antage/eventsource"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"ls/internal/app/plantform_tool/form"
	"ls/internal/pkg/common"
	"ls/internal/pkg/lib/logger"
	"ls/internal/pkg/lib/redis"

	upZip "ls/internal/pkg/lib/zip"
	"net/http"
	"os"
	"time"
)

type File struct {
	common.BaseHandler
}
/*上传*/
func (h File) Upload(c *gin.Context){
	//接收参数
	var form form.UploadFileForm
	if err := h.BindParams(c, &form); err != nil {
		h.HandleError(c, err)
		return
	}

	// 判断本地是否存在该文件
	// 存在 ：判断版本是否为最新 且最新，todo:判断是否解压 没解压解压，返回下载完成
	// 不存在｜不是最新版本则
	// 查看redis内是否正在进行
	// redis内没有正在下载任务->保存到redis(新建task_id,保存映射关系和队列)
	//创建一个sse链接
	es := eventsource.New(
		&eventsource.Settings{
			Timeout: 5 * time.Second,
			CloseOnTimeout: false,
			//空闲链接超时时间
			IdleTimeout: 30 * time.Minute,
		}, func(req *http.Request) [][]byte {
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
			}
		})
	es.SendRetryMessage(3 * time.Second)
	es.ServeHTTP(c.Writer, c.Request)
	var filemap = make(map[string]int)
	go rollback(es,filemap)
}
func rollback(es eventsource.EventSource,filemap map[string]int) {
	//轮询下载文件列表
	for {
		//如果链接断开了停止协程
		if es.ConsumersCount() == 0 {
			es.Close()
			break
		}
		time.Sleep(time.Duration(100))
		var empty = true
		for i,o:=range filemap{
			if o == 0{
				continue
			}
			empty = false
			fmt.Println("文件路径 = ", i, "taskid = ", o)
			re,err:= redis.RedisGetString(o)
			if err != nil || re == ""{
				//todo:没有查询到数据 或报错
			}
			//检查下载情况
			//1、等待
			//2、开始
			//3、完成 区分下载和解压的后续工作
			es.SendEventMessage("over", "message", i)
		}
		if empty {
			break
		}
	}
	es.Close()
}
//下载
func (h File) Download(c *gin.Context){
	//获取请求参数
	/*params,err := h.GetParams(c)*/

}

/*四、查询本地文件*/
func (h File) CheckExists(c *gin.Context){
	//参数解析
	var request form.CheckExistsForm
	if err := h.BindParams(c, &request); err != nil {
		h.HandleError(c, err)
		return
	}
	path := request.Path
	var re = form.CheckExistsRBForm{
		Path:path,
		Exists:false,
	}
	if _, err := os.Stat(path); err != nil{
		re.Exists = true
	}
	h.Success(c,re)
}
/*十、读取文件*/
func (h File) ReadFromFile(c *gin.Context) {
	//参数解析
	var request form.ReadFromFileForm
	if err := h.BindParams(c, &request); err != nil {
		h.HandleError(c, err)
		return
	}
	var re = form.ReadFromFileRBForm{
		ReadFromFileForm:request,
	}
	f, err := os.Open(request.Path)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("ReadFromFile open file err :%s",request.Path),
			zap.String("errMsg",err.Error()))
		h.Success(c,re)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("ReadFromFile close file err :%s",request.Path),
				zap.String("errMsg",err.Error()))
		}
	}(f)

	data,err :=ioutil.ReadAll(f)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("ReadFromFile read file err :%s",request.Path),
			zap.String("errMsg",err.Error()))
		h.Success(c,re)
		return
	}
	re.Data = string(data)
	h.Success(c,re)
	return
}

/*十一、压缩*/
func (h File) Compress(c *gin.Context) {
	//参数解析
	var request form.CompressForm
	if err := h.BindParams(c, &request); err != nil {
		h.HandleError(c, err)
		return
	}
	var re = form.CompressRBForm{Ret: true}
	err := upZip.Compress(request.Source,request.Dest)
	if err != nil{
		re.Ret = false
		logger.Logger.Error("Compress err",zap.String("errMsg",err.Error()))
	}
	h.Success(c,re)
	return
}

/*十二、解压*/
func (h File) Decompress(c *gin.Context) {
	//参数解析
	var request form.DecompressForm
	if err := h.BindParams(c, &request); err != nil {
		h.HandleError(c, err)
		return
	}
	var re = form.DecompressRBForm{Ret: true}
	err := upZip.Decompress(request.Source,request.Dest)
	if err != nil{
		re.Ret = false
		logger.Logger.Error("Decompress err",zap.String("errMsg",err.Error()))
	}
	h.Success(c,re)
	return
	//目标文件夹不存在则创建

	h.Success(c,re)
	return
}

