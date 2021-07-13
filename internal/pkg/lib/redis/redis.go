package redis

import (
	"github.com/gomodule/redigo/redis"
	"ls/internal/app/plantform_tool"
	"ls/internal/app/plantform_tool/config"
	"ls/internal/pkg/lib/logger"
	"time"
)

/*链接池初始化*/
func PoolInitRedis(conf config.RedisConfig) {
	plantform_tool.RedisPool = &redis.Pool{
		MaxIdle:     conf.MaxIdle,//空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   conf.MaxActive,//最大数
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Server)
			if err != nil {
				return nil, err
			}
			if plantform_tool.RedisConfig.Password != "" {
				if _, err := c.Do("AUTH",conf.Password); err != nil {
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

	logger.Logger.Info("redis初始化成功")
}
/*key-value set*/
func RedisSetString(key,value interface{}) (string, error) {
	con := plantform_tool.RedisPool.Get()
	if err := con.Err(); err != nil {
		return "", err
	}
	defer con.Close()
	return redis.String(con.Do("SET", key,value))
}

/*key-value get string*/
func RedisGetString(key interface{}) (string, error) {
	con := plantform_tool.RedisPool.Get()
	if err := con.Err(); err != nil {
		return "", err
	}
	defer con.Close()

	return redis.String(con.Do("GET", key))
}
/*自增的数字*/
func INCRInt64(key interface{}) (int64, error){
	con := plantform_tool.RedisPool.Get()
	if err := con.Err(); err != nil {
		return 0, err
	}
	defer con.Close()
	return redis.Int64(con.Do("INCR", key))
}
/*list插入数据 左进*/
func LPush(key,value interface{}) (int64, error){
	con := plantform_tool.RedisPool.Get()
	if err := con.Err(); err != nil {
		return 0, err
	}
	defer con.Close()
	return redis.Int64(con.Do("LPUSH", key,value))
}
/*list移除数据 右出*/
func RPop(key interface{}) (int, error){
	con := plantform_tool.RedisPool.Get()
	if err := con.Err(); err != nil {
		return 0, err
	}
	defer con.Close()
	return redis.Int(con.Do("RPOP", key))
}
/*list数据长度*/
func LLen(key interface{}) (int, error){
	con := plantform_tool.RedisPool.Get()
	if err := con.Err(); err != nil {
		return 0, err
	}
	defer con.Close()
	return redis.Int(con.Do("LLen", key))
}
/*根据索引取值*/
func LIndex(key,index interface{}) (int, error){
	con := plantform_tool.RedisPool.Get()
	if err := con.Err(); err != nil {
		return 0, err
	}
	defer con.Close()
	return redis.Int(con.Do("LINDEX", key,index))
}
/*根据值移除*/
func LRem(key,value interface{}) (int64, error){
	con := plantform_tool.RedisPool.Get()
	if err := con.Err(); err != nil {
		return 0, err
	}
	defer con.Close()
	return redis.Int64(con.Do("LREM", key,0,value))
}
//执行自定义redis命令
func RedisExec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := plantform_tool.RedisPool.Get()
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