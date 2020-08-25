package main

import (
	"flag"
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

	var poolNum int
	flag.IntVar(&poolNum, "n", 5000, "")
	flag.Parse()

	fmt.Println("goroutine count:", poolNum)

	go Manager.Start()
	var err error
	poller, err = netpoll.New(nil)
	fmt.Println("poller:", poller)
	if err != nil {
		log.Fatal(err)
	}
	pool = NewPool(poolNum, poolNum, 50)

	go IntiWsNotifyConsumer()

	r := setUpRouter()
	r.Run(":8090")
}
