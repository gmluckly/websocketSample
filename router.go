package main

import (
	"github.com/gin-gonic/gin"
)

func setUpRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.GET("/api/websocket/ws/notify", WsNotifyHandler)

	return r
}
