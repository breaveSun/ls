package app

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"
)

func CheckoutPort(port int)string{
	res := ""
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr %d",port)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.Stdout = &outBytes
	cmd.Run()
	resStr := outBytes.String()

	exists := strings.Contains(strings.ToUpper(resStr), "LISTENING")
	//uplog.Debug("port occupy = ", exists)
	if !exists {
		return res
	}

	r := regexp.MustCompile(`\s\d+\s`).FindAllString(resStr, -1)
	if len(r) > 0 {
		res = strings.TrimSpace(r[0])

	}
	uplog.Debug("port in use : ", res)
	return res
}
