package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CodeData struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
// Success 成功
func (BaseHandler) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK,data)
}
// ErrorCode 成功(自定义code)
func (BaseHandler) ErrorCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, CodeMessage{
		Code:    code,
		Message: message,
	})
}
// InvalidParameter 参数无效
func (BaseHandler) InvalidParameter(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Error{
		msg,
	})
}
// ServerError 服务异常
func (BaseHandler) ServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Error{
		msg,
	})
}
// HandleError 自定义错误处理
func (b BaseHandler) HandleError(c *gin.Context, err error) {
		// 内部错误
	switch err {
		//服务名称错误
	case ErrServerName:
		b.ErrorCode(c, 1001, err.Error())
		//其他错误……
	default:
		b.ServerError(c, err.Error())
	}
}
