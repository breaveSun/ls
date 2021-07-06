package common

import "errors"
var (
	ErrServerName    = errors.New("服务名称错误")
)
type Error struct {
	Error string `json:"error"`
}
type CodeMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}