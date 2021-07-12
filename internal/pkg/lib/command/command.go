package command

import (
	"fmt"
	"os/exec"
	"strings"
)

func Test(name string) bool{
	output, _ :=  exec.Command("cmd", "/C", "tasklist").Output()
	n := strings.Index(string(output), "System")
	if n == -1 {
		fmt.Println("windows system abnormal")
		return false
	}
	a := strings.Index(string(output), `\r\n`+name+`\s`)
	if a >=0{
		return true
	}
	return false
}
