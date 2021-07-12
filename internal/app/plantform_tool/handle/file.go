package handle

import "C"
import (
	"errors"
	"fmt"
	"github.com/antage/eventsource"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
	"io/ioutil"
	"ls/internal/app/plantform_tool"
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
/*	var form form.UploadFileForm
	if err := h.BindParams(c, &form); err != nil {
		h.HandleError(c, err)
		return
	}*/
	//获取参数
	result,_:=h.GetParams(c)
	//判断类型
	switch gjson.Get(result,plantform_tool.TransType).String() {
	case plantform_tool.UploadKey:
		//上传要判断本地路径是否存在
		path := gjson.Get(result,plantform_tool.LocalPathKey).String()
		if _, err := os.Stat(path); err != nil{
			h.HandleError(c,errors.New(fmt.Sprintf("not found in path : %s",path)))
			return
		}
		//是否需要压缩
		needZip := gjson.Get(result,plantform_tool.NeedZipKey).Int()
		TargetPath := gjson.Get(result,plantform_tool.ZipTargetKey).String()
		if needZip == plantform_tool.NeedZip && TargetPath != ""{
			//压缩
			err := upZip.Compress(path,TargetPath)
			if err != nil{
				h.HandleError(c,errors.New(fmt.Sprintf("zip err %s -> %s reason:%s",path,TargetPath,err.Error())))
				return
			}
		}
	case plantform_tool.DownLoadKey:
	}

	//生成一个taskId
	taskId,err:=redis.INCRInt64(plantform_tool.TaskIdRDKey)
	if err!=nil{
		h.HandleError(c,errors.New(fmt.Sprintf("taskId create err:%s",err.Error())))
		return
	}
	//插入任务id
	result,err=sjson.Set(result,plantform_tool.TaskIdKey,taskId)
	if err!=nil{
		h.HandleError(c,errors.New(fmt.Sprintf("set %s err:%s",plantform_tool.TaskIdKey,err.Error())))
		return
	}

	//插入创建时间
	result,err=sjson.Set(result,plantform_tool.CreateTimeKey,time.Now())
	if err!=nil{
		h.HandleError(c,errors.New(fmt.Sprintf("set %s err:%s",plantform_tool.CreateTimeKey,err.Error())))
		return
	}
	//把任务保存到redis
	_,err = redis.RedisSetString(taskId,result)
	if err!=nil{
		h.HandleError(c,errors.New(fmt.Sprintf("task %s save redis err:%s",plantform_tool.TaskIdRDKey,err.Error())))
		return
	}
	//把任务id增加到需要回调的列表
	_,err = redis.LPush(plantform_tool.TaskPollingListRDKey,taskId)
	if err != nil {
		h.HandleError(c,errors.New(fmt.Sprintf("task polling list %d save redis err:%s",taskId,err.Error())))
		return
	}
	C.notify()
	h.Success(c,nil)
	return
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
}

func StartFileTransferProgressListen(){
	//获取任务列表长度
	for {
		go FileTransferProgressListen()
		time.Sleep(time.Duration(plantform_tool.ServerConfig.FileTransferListenInterval))
	}

}
func FileTransferProgressListen(){
	//读取当前需要监听的taskid
	taskListLen,err := redis.LLen(plantform_tool.TaskPollingListRDKey)
	if err != nil{
		return
	}

	for i:=0;i<taskListLen;i++{
		taskId,err := redis.LIndex(plantform_tool.TaskPollingListRDKey,i)
		if err != nil{
			continue
		}
		//根据taskId获取任务内容
		task,err := redis.RedisGetString(taskId)
		if err != nil{
			continue
		}
		//判断任务是否完成
		status := gjson.Get(task,plantform_tool.TaskStatusKey).Int()
		switch status {
		case plantform_tool.TaskStatusOver:
			//判断是否需要解压或者压缩
			needZip := gjson.Get(task,plantform_tool.NeedZipKey).Int()
			switch needZip {
			case plantform_tool.NeedUnZip:
				//把本地路径解压到
				localPath := gjson.Get(task,plantform_tool.LocalPathKey).String()
				TargetPath := gjson.Get(task,plantform_tool.ZipTargetKey).String()
				err := upZip.Decompress(localPath,TargetPath)
				if err != nil{
					logger.Logger.Error(fmt.Sprintf("解压缩失败：%s -> %s",localPath,TargetPath))
					continue
				}
				//给服务器发送通知
				//服务器地址
				callBackUrl := gjson.Get(task,plantform_tool.CallBackUrlKey).String()
				//剩余回调次数
				surplusCount := gjson.Get(task,plantform_tool.CallBackTryCountKey).Int()
				if callBackUrl == "" || surplusCount == 0{
					//移除队列
					if _, err = redis.LRem(plantform_tool.TaskPollingListRDKey, taskId);err != nil {
						continue
					}
				} else {
					//需要回调
					fileSize := gjson.Get(task,plantform_tool.FileSizeKey).String()
					callBackData := gjson.Get(task,plantform_tool.CallBackDataKey).String()
					//整理回调参数
					request := form.FileTransferRequestForm{
						FileSize:fileSize,
						CallBackData: callBackData,
					}
					//请求服务器
					re,err:=common.CallServer(callBackUrl,request)
					if err != nil{
						continue
					}
					//解析回调结果
					var response form.ServerResponseForm
					err = jsoniter.UnmarshalFromString(re,&response)
					if err != nil{
						continue
					}

					if response.Code == plantform_tool.ServerResponseCode {
						/*回调成功*/
						//保存次数
						_,_ = sjson.Set(task,plantform_tool.CallBackTryCountKey,surplusCount)
						//移除队列
						if _, err = redis.LRem(plantform_tool.TaskPollingListRDKey, taskId);err != nil {
							continue
						}
					} else {
						/*回调失败*/
						//减少回调次数
						surplusCount := gjson.Get(task,plantform_tool.CallBackTryCountKey).Int()-1
						if surplusCount == 0{
							//保存次数
							_,_ = sjson.Set(task,plantform_tool.CallBackTryCountKey,surplusCount)
							//移除会监控队列
							if _, err = redis.LRem(plantform_tool.TaskPollingListRDKey, taskId);err != nil {
								continue
							}
						} else {
							//保存回调次数
							_,_ = sjson.Set(task,plantform_tool.CallBackTryCountKey,surplusCount)
						}
					}
				}
			}
		}
	}
}

//下载//todo::考虑与upload合并
func (h File) Download(c *gin.Context){
	//获取请求参数
	/*params,err := h.GetParams(c)*/

}
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


