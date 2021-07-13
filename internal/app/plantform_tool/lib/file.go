package lib

import (
	"errors"
	"fmt"
	"github.com/antage/eventsource"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
	"ls/internal/app/plantform_tool"
	"ls/internal/app/plantform_tool/constant"
	"ls/internal/app/plantform_tool/form"
	"ls/internal/pkg/common"
	"ls/internal/pkg/lib/logger"
	"ls/internal/pkg/lib/redis"
	upZip "ls/internal/pkg/lib/zip"
	"os"
	"time"
)

// UploadBefore /*上传任务前置操作*/
func UploadBefore(request string) error{
	//上传要判断本地路径是否存在
	LocalPath := gjson.Get(request,constant.LocalPathKey).String()
	if LocalPath == ""{
		return errors.New(fmt.Sprintf("%s can't empty",constant.LocalPathKey))
	}
	//是否需要压缩
	needZip := gjson.Get(request,constant.NeedZipKey).Int()
	TargetPath := gjson.Get(request,constant.ZipTargetKey).String()
	if needZip == constant.NeedZip{
		//需要压缩
		if TargetPath == ""{
			return errors.New(fmt.Sprintf("needZip the %s is must",constant.LocalPathKey))
		}
		if _, err := os.Stat(TargetPath); err != nil{
			return errors.New(fmt.Sprintf("not found in path : %s",TargetPath))
		}
		err := upZip.Compress(TargetPath,LocalPath)
		if err != nil{
			return errors.New(fmt.Sprintf("zip err %s -> %s reason:%s",TargetPath,LocalPath,err.Error()))
		}
	} else {
		// 不需要压缩，上传前要判断本地路径是否存在
		if _, err := os.Stat(LocalPath); err != nil{
			return errors.New(fmt.Sprintf("not found in path : %s",LocalPath))
		}
	}
	return nil
}

// StartFileTransferProgressListen /*开启文件上传任务进度监控*/
func StartFileTransferProgressListen(){
	for {
		//开始监听未完成上传下载任务进度
		go fileTransferProgressListen()
		//一定的时间间隔轮询
		time.Sleep(time.Duration(plantform_tool.ServerConfig.FileTransferListenInterval))
	}
}

/*监听未完成上传下载任务进度*/
func fileTransferProgressListen(){
	//读取当前需要监听的taskid
	taskListLen,err := redis.LLen(constant.TaskPollingListRDKey)
	if err != nil{
		return
	}

	for i:=0;i<taskListLen;i++{
		taskId,err := redis.LIndex(constant.TaskPollingListRDKey,i)
		if err != nil{
			continue
		}
		//根据taskId获取任务内容
		task,err := redis.RedisGetString(taskId)
		if err != nil{
			continue
		}
		//判断任务是否完成
		status := gjson.Get(task,constant.TaskStatusKey).Int()
		switch status {
		case constant.TaskStatusOver:
			err:=downLoadAfter(task,taskId)
			if err !=nil{
				logger.Logger.Error("FileTransferProgressListen",zap.String("errMsg",err.Error()))
			}
		}
	}
}

/*下载任务完成之后*/
func downLoadAfter(task string,taskId int)error{
	//判断是否需要解压或者压缩
	needZip := gjson.Get(task,constant.NeedZipKey).Int()
	switch needZip {
	//需要解压
	case constant.NeedUnZip:
		//把下载之后的本地路径解压到
		localPath := gjson.Get(task, constant.LocalPathKey).String()
		TargetPath := gjson.Get(task, constant.ZipTargetKey).String()
		err := upZip.Decompress(localPath, TargetPath)
		if err != nil {
			return errors.New(fmt.Sprintf("解压缩失败：%s -> %s", localPath, TargetPath))
		}
	}

	//判断是否给服务器发送通知（重试次数和回调地址）
	callBackUrl := gjson.Get(task,constant.CallBackUrlKey).String()//服务器地址
	surplusCount := gjson.Get(task,constant.CallBackSurplusCountKey).Int()//剩余回调次数
	if callBackUrl == "" || surplusCount == 0{
		//移除队列
		if _, err := redis.LRem(constant.TaskPollingListRDKey, taskId);err != nil {
			return errors.New(fmt.Sprintf("把已完成任务taskId=%d 移除监控队列失败 : %s",taskId,err.Error()))
		}
		return nil
	}

	//回调服务器上传|下载进度
	fileSize := gjson.Get(task,constant.FileSizeKey).String()
	callBackData := gjson.Get(task,constant.CallBackDataKey).String()
	//整理回调参数
	request := form.FileTransferRequestForm{
		FileSize:fileSize,
		CallBackData: callBackData,
	}
	//请求服务器
	re,err:=common.CallServer(callBackUrl,request)
	if err != nil{
		return errors.New(fmt.Sprintf("上传下载结果回调失败 taskId = %d url = %s err = %s",taskId,callBackUrl,err.Error()))
	}
	//解析回调结果
	var response form.ServerResponseForm
	err = jsoniter.UnmarshalFromString(re,&response)
	if err != nil{
		return errors.New(fmt.Sprintf("上传下载结果回调参数解析失败 taskId = %d url = %s response = %s err = %s",taskId,callBackUrl,re,err.Error()))
	}
	if surplusCount > 0 {
		surplusCount-- //减少回调次数
		_,_ = sjson.Set(task,constant.CallBackSurplusCountKey,surplusCount)//保存次数
	}
	/*回调成功|剩余回调次数=0*/
	if response.Code == plantform_tool.ServerResponseCode || surplusCount == 0{
		if _, err = redis.LRem(constant.TaskPollingListRDKey, taskId);err != nil {//移除队列
			return errors.New(fmt.Sprintf("把已完成任务 taskId=%d （回调成功|回调次数=0）移除监控队列失败 : %s",taskId,err.Error()))
		}
	}
	return nil
}

// Rollback /*SSE发送数据给客户端*/
func Rollback(es eventsource.EventSource,filemap map[string]int) {
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


