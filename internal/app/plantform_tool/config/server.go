package config

type ServerConfig struct {
	Restart int
	Port int
	//监听文件上传和下载轮询时间间隔
	FileTransferListenInterval int
}

