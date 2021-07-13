package form

// ExecCommandFrom /*六、CMD 执行*/
type ExecCommandFrom struct {
	Command string `json:"Command"`
	Args []string `json:"Args"`
}
type ExecCommandRBFrom struct {
	ErrorInfo string `json:"ErrorInfo"`
}

// KillProcessFrom /*八、杀进程*/
type KillProcessFrom struct {
	AppName string `json:"AppName"`
}
type KillProcessRBFrom struct {
	KillProcessFrom
	Ret bool `json:"Ret"`
}