package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)
var (
	Validate   *validator.Validate
	Translator ut.Translator
)
type BaseHandler struct{}
// BindParams 绑定参数到结构体
func (BaseHandler) BindParams(c *gin.Context, d interface{}) error {
	if err := c.ShouldBind(d); err != nil {
		// Invalid params
		return err
	}
	// Translate error message
	err := Validate.Struct(d)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			var sliceErrs []string
			for _, err := range err.(validator.ValidationErrors) {
				sliceErrs = append(sliceErrs, err.Translate(Translator))
			}
			return errors.New(strings.Join(sliceErrs, ","))
		}
		return err
	}
	return nil
}
