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
