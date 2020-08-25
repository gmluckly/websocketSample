package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"strconv"
)

func WsNotifyHandler(c *gin.Context) {
	fmt.Println("ws.......")

	userId, check := websocketCheckSession(c)
	if !check {
		fmt.Println("websocket init error ")
		return
	}
	w := c.Writer
	r := c.Request
	conn, _, _, err := ws.UpgradeHTTP(r, w)

	if err != nil {
		fmt.Println("Failed to set websocket upgrade:", err)
		return
	}
	client := &Client{Id: userId, Socket: conn}
	fmt.Println("start web socket,send to Manager")
	Manager.Register <- client
	fmt.Println("websocke init success")

	desc := netpoll.Must(netpoll.HandleRead(conn))
	//fmt.Println("desc:", desc)
	poller.Start(desc, func(ev netpoll.Event) {
		if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
			//fmt.Println("websocket stop connect:", ev)
			Manager.Unregister <- client
			poller.Stop(desc)
			return
		}
		pool.Schedule(func() {
			client.Read()
		})
	})
}

func websocketCheckSession(c *gin.Context) (int64, bool) {
	userIdStr, _ := c.GetQuery("uid")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	return userId, true
}
