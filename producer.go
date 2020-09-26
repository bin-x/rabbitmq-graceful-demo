package main

import (
	"github.com/bin-x/rabbitmq-graceful-demo/lib"
	"log"
)

func main() {
	var mq lib.MqClient
	err := mq.Connect("amqp://user:password@ip:port/yourhost")
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()

	xName := "testX"
	if err := mq.Publish([]byte("this is message"), xName, "fanout", ""); err != nil {
		log.Println(err)
	} else {
		log.Println("send message: this is message")
	}
}
