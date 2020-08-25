package main

import (
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"net"
	"sync"
	//"github.com/gobwas/ws"
	//"github.com/gorilla/websocket"
)

type msg struct {
	EventType string      `json:"eventType"`
	Data      interface{} `json:"data"`
}

var resp []byte

func init() {
	var pingMsg = &msg{EventType: "pong", Data: ""}
	resp, _ = json.Marshal(pingMsg)
}

type Client struct {
	Id     int64
	Socket net.Conn
}

type ClientManager struct {
	mutex      *sync.RWMutex
	Clients    map[int64]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

var Manager = ClientManager{
	mutex:      new(sync.RWMutex),
	Register:   make(chan *Client, 10),
	Unregister: make(chan *Client, 10),
	Clients:    make(map[int64]*Client, 1000),
}

func (Manager *ClientManager) Start() {
	for {
		select {
		case client := <-Manager.Register:
			Manager.mutex.Lock()
			Manager.Clients[client.Id] = client
			Manager.mutex.Unlock()
		case client := <-Manager.Unregister:
			//fmt.Println("delete fro map:", client.Id)
			if _, ok := Manager.Clients[client.Id]; ok {
				Manager.mutex.Lock()
				delete(Manager.Clients, client.Id)
				Manager.mutex.Unlock()
			}
		}
	}
}

func (c *Client) Read() {
	var (
		r       = wsutil.NewReader(c.Socket, ws.StateServerSide)
		w       = wsutil.NewWriter(c.Socket, ws.StateServerSide, ws.OpText)
		decoder = json.NewDecoder(r)
		//encoder = json.NewEncoder(w)
	)

	hdr, err := r.NextFrame()
	if err != nil {
		fmt.Println("000", err)
	}

	if hdr.OpCode == ws.OpClose {
		fmt.Println("111", io.EOF)
	}

	var req msg
	if err := decoder.Decode(&req); err != nil {
		fmt.Println("222", err)
	}
	//fmt.Println("req:", req)

	if req.EventType == "ping" {
		//resp := msg{EventType: "pong", Data: ""}
		//if err := encoder.Encode(&resp); err != nil {
		//	fmt.Println("333", err)
		//}
		w.Write(resp)
		if err = w.Flush(); err != nil {
			fmt.Println("444", err)
		}
	} else {
		fmt.Println("error can not get other msg expect ping", err)
	}
}

/*
func (c *Client) Read() {
	for {
		defer c.Socket.Close()
		_, message, err := c.Socket.ReadMessage()
		fmt.Println("77  message:", message)
		if err != nil {
			fmt.Println("80  read message err:", err)
			Manager.Unregister <- c
			//c.Socket.Close()
			break
		}
		if strings.EqualFold(string(message), "ping") {
			//fmt.Println("get message from client:", message)
			c.Socket.WriteMessage(websocket.TextMessage, []byte("pong"))
		} //else {
		//	sender := strconv.FormatInt(c.Id, 10)
		//	jsonMessage, _ := json.Marshal(&Message{Sender: sender, Content: string(message)})
		//	c.Socket.WriteMessage(websocket.TextMessage, jsonMessage)
		//}
	}
}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func MessageFromMQ(userId int64, message []byte) {
	if client, ok := Manager.Clients[userId]; ok {
		sender := strconv.FormatInt(userId, 10)
		jsonMessage, _ := json.Marshal(&Message{Sender: sender, Content: string(message)})
		client.Socket.WriteMessage(websocket.TextMessage, jsonMessage)
	}

}
*/

func sendMQMsg(userId int64, eventType string, message interface{}) {
	if client, ok := Manager.Clients[userId]; ok {
		w := wsutil.NewWriter(client.Socket, ws.StateServerSide, ws.OpText)
		msg := make(map[string]interface{})
		msg["eventType"] = eventType
		msg["data"] = message
		byteMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("0000", err)
			return
		}
		w.Write(byteMsg)
		if err := w.Flush(); err != nil {
			fmt.Println("000", err)
		}
	}
}
