package handle
/*
#include <SemaphoreForGo.h>
*/
import "C"
import "C"
import (
	"errors"
	"fmt"
	"github.com/antage/eventsource"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
	"io/ioutil"
	"ls/internal/app/plantform_tool/constant"
	"ls/internal/app/plantform_tool/form"
	"ls/internal/app/plantform_tool/lib"
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

// UploadDownLoad /*上传|下载*/
func (h File) UploadDownLoad(c *gin.Context){
	//获取参数
	result,_:=h.GetParams(c)
	//判断类型
	switch gjson.Get(result,constant.TransType).String() {
	case constant.UploadKey:
		//上传任务前置操作
		err := lib.UploadBefore(result)
		h.HandleError(c,err)
		return
	case constant.DownLoadKey:
		//下载
	}

	//生成一个taskId
	taskId,err:=redis.INCRInt64(constant.TaskIdRDKey)
	if err!=nil{
		h.HandleError(c,errors.New(fmt.Sprintf("taskId create err:%s",err.Error())))
		return
	}
	//插入任务id
	result,err=sjson.Set(result,constant.TaskIdKey,taskId)
	if err!=nil{
		h.HandleError(c,errors.New(fmt.Sprintf("set %s err:%s",constant.TaskIdKey,err.Error())))
		return
	}

	//插入创建时间
	result,err=sjson.Set(result,constant.CreateTimeKey,time.Now())
	if err!=nil{
		h.HandleError(c,errors.New(fmt.Sprintf("set %s err:%s",constant.CreateTimeKey,err.Error())))
		return
	}

	//插入剩余回调次数
	maxTry:=gjson.Get(result,constant.MaxTryKey).Int()
	result,err=sjson.Set(result,constant.CallBackSurplusCountKey,maxTry)
	if err!=nil{
		h.HandleError(c,errors.New(fmt.Sprintf("set %s err:%s",constant.CallBackSurplusCountKey,err.Error())))
		return
	}

	//把任务保存到redis
	_,err = redis.RedisSetString(taskId,result)
	if err!=nil{
		h.HandleError(c,errors.New(fmt.Sprintf("task %s save redis err:%s",constant.TaskIdRDKey,err.Error())))
		return
	}

	//把任务id增加到需要回调的列表
	_,err = redis.LPush(constant.TaskPollingListRDKey,taskId)
	if err != nil {
		h.HandleError(c,errors.New(fmt.Sprintf("task polling list %d save redis err:%s",taskId,err.Error())))
		return
	}

	//给文件服务发送信号
	C.notify()
	h.Success(c,nil)
	return
}

// CheckExists /*四、查询本地文件*/
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

// ReadFromFile /*十、读取文件*/
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

// Compress /*十一、压缩*/
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

// Decompress /*十二、解压*/
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
}


/*
//下载//todo::考虑与upload合并
func (h File) Download(c *gin.Context){
	//获取请求参数
	//params,err := h.GetParams(c)

}*/

func (h File) Js(c *gin.Context){
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
	go lib.Rollback(es,filemap)
}



