package handle

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sys/windows/registry"
	"ls/internal/app/plantform_tool/form"
	"ls/internal/pkg/common"
	upcommand "ls/internal/pkg/lib/command"
	"ls/internal/pkg/lib/logger"
)

const LocalMachine = "HKEY_LOCAL_MACHINE"
const CurrentUser = "HKEY_CURRENT_USER"
type WindowsServer struct {
	common.BaseHandler
}
//五、读取注册表
func (h WindowsServer) ReadRegistry(c *gin.Context) {
	var request  []form.ReadRegistryFrom
	if err := h.BindParams(c, &request); err != nil {
		h.HandleError(c, err)
		return
	}
	var registInfoList []form.ReadRegistryRBFrom
	for _,o:=range request{
		var re = form.ReadRegistryRBFrom{
			ReadRegistryFrom:o,
		}
		var rootKey registry.Key
		switch o.Root {
		case LocalMachine:
			rootKey = registry.CURRENT_USER
		case CurrentUser:
			rootKey = registry.CURRENT_USER
		}
		k, err := registry.OpenKey(rootKey, o.Path, registry.ALL_ACCESS)
		if err != nil {
			logger.Logger.Error("get "+o.Root+" err",zap.String("errMsg",err.Error()))
		}
		re.Value,_,err=k.GetStringValue(o.Key)
		if err != nil{
			logger.Logger.Error("get " + o.Root + "'s " + o.Key + " error ",zap.String("errMsg",err.Error()))
		}
		registInfoList = append(registInfoList,re)
	}

	h.Success(c,registInfoList)
}


/*七-1、App 运行状态检测（一次性）*/
func (h WindowsServer) CheckRunning(c *gin.Context) {
	var request form.CheckRunningFrom
	if err := h.BindParams(c, &request); err != nil {
		h.HandleError(c, err)
		return
	}
	if request.AppName !=""{
		var re  = form.CheckRunningAppNameRBFrom{
			AppName: request.AppName,
			Running: upcommand.Test(request.AppName),
		}
		h.Success(c,re)
		return
	} else if request.MemName != ""{
		var re  = form.CheckRunningAppNameRBFrom{
			AppName: request.AppName,
			Running: upcommand.Test(request.AppName),
		}
		h.Success(c,re)
		return
	}
	h.HandleError(c,errors.New("param err"))
	return
}
/*七-2、App 运行状态检测（持续检测）*/
func (h WindowsServer) RunningStatus(c *gin.Context) {
	var request form.RunningStatusFrom
	if err := h.BindParams(c, &request); err != nil {
		h.HandleError(c, err)
		return
	}
}