package main

import (
	"fmt"
	"github.com/mailru/easygo/netpoll"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

var poller netpoll.Poller
var pool *Pool

func main() {
	fmt.Println("start the websocket server...")
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		log.Println(http.ListenAndServe(":10000", nil))
	}()

	go Manager.Start()
	var err error
	poller, err = netpoll.New(nil)
	fmt.Println("poller:", poller)
	if err != nil {
		log.Fatal(err)
	}
	poolNum := 100
	pool = NewPool(poolNum, poolNum, 50)

	go IntiWsNotifyConsumer()

	r := setUpRouter()
	r.Run(":8090")
}
