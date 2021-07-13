package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ls/internal/app/plantform_tool/form"
	"ls/internal/pkg/common"
	"ls/internal/pkg/lib/logger"
	"os/exec"
	"strings"
)

type Command struct {
	common.BaseHandler
}

// ExecCommand /*六、CMD 执行*/
func (h Command) ExecCommand(c *gin.Context){
	var request  form.ExecCommandFrom
	if err := h.BindParams(c, &request); err != nil {
		h.HandleError(c, err)
		return
	}
	var re form.ExecCommandRBFrom
	var cmd = fmt.Sprintf("%s %s", request.Command,strings.Join(request.Args," "))
	if _, err := exec.Command("cmd", "/C",cmd).
		Output();err != nil{
			re.ErrorInfo = err.Error()
			logger.Logger.Error("ExecCommand err",zap.String(cmd,re.ErrorInfo))
	}
	h.Success(c,re)
	return
}

// KillProcess /*八、杀进程*/
func (h Command) KillProcess(c *gin.Context) {
	var request form.KillProcessFrom
	if err := h.BindParams(c, &request); err != nil {
		h.HandleError(c, err)
		return
	}
	var re = form.KillProcessRBFrom{
		KillProcessFrom:request,
	}
	if _, err := exec.Command("cmd", "/C", `taskkill /f /t /im ` +
		request.AppName).Output();err != nil{
		logger.Logger.Error(fmt.Sprintf("KillProcess %s err",request.AppName),zap.String("errMsg",err.Error()))
	} else {
		re.Ret = true
	}
	h.Success(c,re)
	return
}