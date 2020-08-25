package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	//"ilikereader.com/link_target_pool_system/websocket"
	"log"
)

func IntiWsNotifyConsumer() {
	var wsNotifyConn *amqp.Connection
	var wsNotifyChannel *amqp.Channel

	queue := "ws_notify.queue.100"
	fmt.Println("init ws_mq,queueName:", queue)

	host := "172.16.3.131"
	url := "amqp://guest:guest@" + host
	var err error
	wsNotifyConn, err = amqp.Dial(url)
	failOnErr(err, "connect ws nofity queue fail")
	wsNotifyChannel, err = wsNotifyConn.Channel()
	failOnErr(err, "channel error")
	defer wsNotifyConn.Close()
	_, err = wsNotifyChannel.QueueDeclare(queue, true, false, false, false, nil)
	failOnErr(err, "declare error ")

	fmt.Println("22222")
	msgs, err := wsNotifyChannel.Consume(queue, "", true, false, false, false, nil)
	failOnErr(err, "get msgs err")
	forever := make(chan bool)
	//userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	go func() {
		for d := range msgs {
			body := d.Body
			type request struct {
				UserId    int64       `json:"userId"`
				EventType string      `json:"eventType"`
				Data      interface{} `json:"data"`
			}
			var req request
			err := json.Unmarshal(body, &req)
			if err != nil {
				fmt.Println("json unmarshal error:", err)
			} else {
				fmt.Println("userId:", req.UserId, ",body:", string(body))

				sendMQMsg(req.UserId, req.EventType, req.Data)
			}
		}
	}()
	<-forever
}

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatal("msg:%s,err:%s", msg, err)
		panic(err)
	}
}
