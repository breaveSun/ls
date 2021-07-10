package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"io/ioutil"
	"ls/internal/pkg/lib/logger"
	"strings"
)

type BaseHandler struct{}
// BindParams 绑定参数到结构体
func (BaseHandler) BindParams(c *gin.Context, d interface{}) error {
	if err := c.ShouldBind(d); err != nil {
		fmt.Println(err)
		// Invalid params
		return err
	}
	// Translate error message
	/*err := Validate.Struct(d)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			var sliceErrs []string
			for _, err := range err.(validator.ValidationErrors) {
				sliceErrs = append(sliceErrs, err.Translate(Translator))
			}
			return errors.New(strings.Join(sliceErrs, ","))
		}
		return err
	}*/
	return nil
}

func (BaseHandler) GetParams(c *gin.Context)(gjson.Result,error){
	var defReturn gjson.Result
	bodyByte, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		logger.Logger.Error("接口解析参数失败",zap.String("errMsg",err.Error()))
		return defReturn,err
	}
	body := string(bodyByte)
	if !gjson.Valid(body) {
		logger.Logger.Error("接口参数不是json格式",zap.String("body：",body))
		return defReturn,err
	}

	result := gjson.Parse(body)
	return result,nil
}
//获取文件路径
func (BaseHandler) GetDirPath(path string) string {
	e :=strings.LastIndex(path, `\`)
	return subString(path, 0, e)
}
//字符串截取
func subString(str string, start, end int) string {
	//rs := []rune(str)
	length := len(str)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return str[start:end]
}
