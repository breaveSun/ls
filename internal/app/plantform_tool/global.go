package plantform_tool

import (
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"ls/internal/app/plantform_tool/config"
)

const (
	// ServerResponseCode /*server response*/
	ServerResponseCode = 200
)


var (
	Config *viper.Viper
	ServerConfig config.ServerConfig
	RedisConfig config.RedisConfig
	RedisPool *redis.Pool
)