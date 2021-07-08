package handle

import (
	"fmt"
	"github.com/antage/eventsource"
	"github.com/gin-gonic/gin"
	"ls/internal/pkg/lib/redis"
	"net/http"
	"time"
)

type File struct {
	//common.BaseHandler
}
//上传
func (h File) Upload(c *gin.Context){
	//接收参数
	/*var form form2.UploadFileForm
	if err := h.BindParams(c, &form); err != nil {
		h.HandleError(c, err)
		return
	}*/

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
	//待定
	/*用于测试sse
	es := eventsource.New(
		&eventsource.Settings{
			Timeout:        2 * time.Second,
			CloseOnTimeout: true,
			IdleTimeout:    2 * time.Second,
			Gzip:           true,
		},
		func(req *http.Request) [][]byte {
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
			}
		},
	)

	es.ServeHTTP(c.Writer, c.Request)

	go func() {
		var id int
		for {
			id++
			time.Sleep(1 * time.Second)
			es.SendEventMessage("blabla", "message", strconv.Itoa(id))
			if es.ConsumersCount() == 0{
				fmt.Println("客户端停止接收")
				break
			}
		}
	}()*/
}