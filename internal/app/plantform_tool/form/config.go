package form

type ServerConfig struct {
	Restart int
	Port int
}

type RedisConfig struct {
	Server string
	Password string
	MaxIdle int //空闲数
	MaxActive int //最大数
}
