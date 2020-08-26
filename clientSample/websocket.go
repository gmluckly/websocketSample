package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"math/rand"
	"time"
)

func ConnectWebSocket(addr string) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		fmt.Println("connect websocket error :", err)
		time.Sleep(20 * time.Second)
		ConnectWebSocket(addr)
		return
	}
	//conn.SetReadLimit()
	conn.SetReadLimit(200000)

	msgChan := make(chan []byte)
	rand.Seed(time.Now().UnixNano())
	s := rand.Intn(50)

	go hearBeat(conn, s)
	go readMessage(addr, conn, msgChan)
	go handler(conn, msgChan)
}

func readMessage(addr string, conn *websocket.Conn, msgChan chan []byte) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("get messge error :", err)
			time.Sleep(20 * time.Second)
			ConnectWebSocket(addr)
			break
		}
		if err == io.EOF {
			continue
		}
		fmt.Println("get message:", string(msg))
		msgChan <- msg
	}
}

type messge struct {
	EventType string `json:"eventType"`
	Data      string `json:"data"`
}

func handler(conn *websocket.Conn, messageChan chan []byte) {
	for {
		select {
		case msg := <-messageChan:
			//TODO send message to chain-gateway and get result to back
			var r messge
			err := json.Unmarshal(msg, &r)
			//fmt.Println("Unmarshal to json", r, " err:", err)
			if err != nil {
				fmt.Println("json.Unmarshal error:", err)
			}
			if r.EventType == "pong" {
				//fmt.Println("get pong from server")
				go hearBeat(conn, 0)
			} else {
				fmt.Println("client get the msg from from mq,conn:", &conn)

			}
		}
	}
}

func hearBeat(c *websocket.Conn, second int) {
	if c == nil {
		return
	}
	//第一次连接随机，以后每隔50s发一次心跳包
	if second == 0 {
		time.Sleep(time.Second * 50)
	} else {
		duration := time.Second * time.Duration(second)
		time.Sleep(duration)
	}
	//fmt.Println("heartBeat ....")
	body := make(map[string]interface{})
	body["eventType"] = "ping"
	body["data"] = ""
	bytesData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("data transaction to json err:", err)
	}
	c.WriteMessage()
	c.WriteMessage(1, bytesData)
}
