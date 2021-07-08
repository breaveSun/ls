package lib

import (
	"github.com/spf13/viper"
	"ls/internal/app/plantform_tool"
)

func InitConf(){
	v := viper.New()
	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		//todo:记录日志
		panic(err)
	}
	plantform_tool.ServerConfig.Restart = v.GetInt("server_restart")
	plantform_tool.ServerConfig.Port = v.GetInt("server_port")
	plantform_tool.RedisConfig.Server = v.GetString("redis_server")
	plantform_tool.RedisConfig.Password = v.GetString("redis_password")
	plantform_tool.RedisConfig.MaxIdle = v.GetInt("max_idle")
	plantform_tool.RedisConfig.MaxActive = v.GetInt("max_active")
	plantform_tool.Config = v
}