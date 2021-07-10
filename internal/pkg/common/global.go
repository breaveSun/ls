package common

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	Validate   *validator.Validate
	Translator ut.Translator
)
