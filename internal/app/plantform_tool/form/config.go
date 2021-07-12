package form

type ServerConfig struct {
	Restart int
	Port int
	//监听文件上传和下载轮询时间间隔
	FileTransferListenInterval int
}

type RedisConfig struct {
	Server string
	Password string
	MaxIdle int //空闲数
	MaxActive int //最大数
}
