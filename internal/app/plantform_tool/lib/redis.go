package lib

import (
	"github.com/gomodule/redigo/redis"
	"ls/internal/app/plantform_tool"
	"ls/internal/app/plantform_tool/clients"
	"time"
)

//链接池初始化
func PoolInitRedis() {
	clients.RedisPool = &redis.Pool{
		MaxIdle:     plantform_tool.RedisConfig.MaxIdle,//空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   plantform_tool.RedisConfig.MaxActive,//最大数
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", plantform_tool.RedisConfig.Server)
			if err != nil {
				return nil, err
			}
			if plantform_tool.RedisConfig.Password != "" {
				if _, err := c.Do("AUTH", plantform_tool.RedisConfig.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
func Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	pool := redis.Pool
	con := pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}