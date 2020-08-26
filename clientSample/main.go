package main

import (
	"flag"
	"fmt"
	"strconv"
)

func main() {
	forever := make(chan bool)
	var count int
	flag.IntVar(&count, "n", 5000, "")
	flag.Parse()

	fmt.Println("goroutine count:", count)
	for i := 0; i < count; i++ {
		ids := strconv.Itoa(i)
		url := "ws://127.0.0.1:8090/api/websocket/ws/notify?uid=" + ids
		go ConnectWebSocket(url)
		//time.Sleep(time.Second * 1)
	}

	<-forever
}
