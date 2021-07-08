package app

import (
	"fmt"
	"ls/internal/pkg/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

/*func CheckoutPort(port int)string{
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


}*/
func SignGrab(){
	logger.Logger.Info("系统信号检测开始")
	quit := make(chan os.Signal)
	signal.Notify(quit)
	s :=<-quit
	fmt.Println(s)
	switch s {
	case syscall.SIGHUP:
		logger.Logger.Fatal("syscall.SIGHUP")
	case syscall.SIGINT:
		logger.Logger.Error("syscall.SIGINT")
	case syscall.SIGQUIT:
		logger.Logger.Fatal("syscall.SIGQUIT")
	case syscall.SIGILL:
		logger.Logger.Fatal("syscall.SIGILL")
	case syscall.SIGTRAP:
		logger.Logger.Fatal("syscall.SIGTRAP")
	case syscall.SIGABRT:
		logger.Logger.Fatal("syscall.SIGABRT")
	case syscall.SIGBUS:
		logger.Logger.Fatal("syscall.SIGBUS")
	case syscall.SIGFPE:
		logger.Logger.Fatal("syscall.SIGFPE")
	case syscall.SIGKILL:
		logger.Logger.Fatal("syscall.SIGKILL")
	case syscall.SIGSEGV:
		logger.Logger.Fatal("syscall.SIGSEGV")
	case syscall.SIGPIPE:
		logger.Logger.Fatal("syscall.SIGPIPE")
	case syscall.SIGALRM:
		logger.Logger.Fatal("syscall.SIGALRM")
	case syscall.SIGTERM:
		logger.Logger.Fatal("syscall.SIGTERM")
	case os.Kill:
		logger.Logger.Fatal("os.Kill")
	}
	fmt.Println("系统信号通知")
	logger.Logger.Error("系统信号通知")
}
