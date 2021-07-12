package common

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"io/ioutil"
	"ls/internal/pkg/lib/logger"
	"net/http"
	"strings"
	"time"
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

func (BaseHandler) GetParams(c *gin.Context)(string,error){
	var defReturn string
	bodyByte, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Logger.Error("接口解析参数失败",zap.String("errMsg",err.Error()))
		return defReturn,err
	}
	body := string(bodyByte)
	if !gjson.Valid(body) {
		logger.Logger.Error("接口参数不是json格式",zap.String("body：",body))
		return defReturn,err
	}

	//result := gjson.Parse(body)
	return body,nil
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
//发送post请求
func CallServer(url string,param interface{})(string,error){
	client := &http.Client{Timeout: 5 * time.Second}
	requestByte,_:=jsoniter.Marshal(param)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestByte))
	defer resp.Body.Close()
	if err != nil {
		return "",err
	}
	result, _ := ioutil.ReadAll(resp.Body)
	return string(result),nil
}