package plantform_tool

import (
	"github.com/spf13/viper"
	"ls/internal/app/plantform_tool/form"
)

var (
	Config *viper.Viper
	ServerConfig form.ServerConfig
	RedisConfig form.RedisConfig
)