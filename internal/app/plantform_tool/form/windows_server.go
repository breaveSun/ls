package form
/*五、读取注册表*/
type ReadRegistryFrom struct {
	Root string `json:"Root" validate:"required"`
	Path string `json:"Path" validate:"required"`
	Key string `json:"Key" validate:"required"`
}

type ReadRegistryRBFrom struct {
	ReadRegistryFrom
	Value string `json:"Value"`
}


/*七-1、App 运行状态检测（一次性）*/
type CheckRunningFrom struct {
	AppName string `json:"AppName"`
	MemName string `json:"MemName"`
}

type CheckRunningAppNameRBFrom struct {
	AppName string `json:"AppName"`
	Running bool `json:"Running"`
}
type CheckRunningMemNameRBFrom struct {
	MemName string `json:"MemName"`
	Running bool `json:"Running"`
}

/*七-2、App 运行状态检测（持续检测）*/
type RunningStatusFrom struct {
	AppName string `json:"AppName"`
	MemName string `json:"MemName"`
}

type RunningStatusAppNameRBFrom struct {
	AppName string `json:"AppName"`
	Running bool `json:"Running"`
}
type RunningStatusMemNameRBFrom struct {
	MemName string `json:"MemName"`
	Running bool `json:"Running"`
}
