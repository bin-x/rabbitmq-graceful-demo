package main

import (
	"github.com/bin-x/rabbitmq-graceful-demo/lib"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {
	var mq lib.MqClient
	err := mq.Connect("amqp://user:password@ip:port/yourhost")
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()

	xName := "testX"
	qName := "testQ"
	cName := "testConsume"
	consume, err := mq.Consume(xName, "fanout", cName, qName, "")
	if err != nil {
		log.Fatalln("error:", err.Error())
	}
	log.Println("start consumer...")

	server := lib.NewServer(work)
	server.SetGraceful(time.Second * 20)
	server.Run(consume)
}

func work(delivery amqp.Delivery) {
	for i := 0; i < 10; i++ {
		log.Println("i:", i)
		time.Sleep(time.Second)
	}
	log.Println("mq's data:", string(delivery.Body[:]))

	//we close the auto-ack, so we need tell the mq server that this message has been processed
	delivery.Ack(false)
}
