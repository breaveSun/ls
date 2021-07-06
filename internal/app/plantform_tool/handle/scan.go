package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ls/internal/app/plantform_tool/clients"
)
type Scan struct {

}
func (h Scan)ScanStart(c *gin.Context){
	err:=clients.RedisClient.Set(c,"scan_status","on",0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ScanStart")
}
func (h Scan)ScanOver(c *gin.Context){
	err:=clients.RedisClient.Set(c,"scan_status","off",0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ScanOver")
}
