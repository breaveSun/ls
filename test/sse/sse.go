package main

import (
	"fmt"
	"github.com/antage/eventsource"
	"net/http"
	"strconv"
	"time"
)

func ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("serveHttp……")
	es := eventsource.New(
		&eventsource.Settings{
			Timeout:        2 * time.Second,
			CloseOnTimeout: true,
			IdleTimeout:    2 * time.Second,
			Gzip:           true,
		},
		func(req *http.Request) [][]byte {
			fmt.Println("set cors..")
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
			}
		},
	)

	es.ServeHTTP(resp, req)

	go func() {
		id := 0
		for {
			fmt.Println("progress = ", id, "count = ", es.ConsumersCount())
			if id<100 && es.ConsumersCount()>0 {
				//es.Close()
				//fmt.Println("progress = ", id, "count = ", es.ConsumersCount())
				es.SendEventMessage("blahblah", "message", strconv.Itoa(id))
				id++
				time.Sleep(100 * time.Millisecond)
			}else if id == 100{
				fmt.Println("close es")
				es.Close()
				return
			}
		}
	}()

}

func main() {
	fmt.Println("running")
	http.HandleFunc("/Test", ServeHTTP)
	err := http.ListenAndServe(":17788", nil)
	fmt.Println(err)
}