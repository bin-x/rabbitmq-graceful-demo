# rabbitmq-graceful-demo
a golang demo to show how to make rabbitmq consumer graceful shutdown

### test
download the code. 


```
$ cd rabbitmq-graceful-demo
$ go get github.com/streadway/amqp
```
open `consumer.go` and `producer.go`, change the rabbitmq's address to yours:

`err := mq.Connect("amqp://user:password@ip:port/yourhost")`

start consumer.go
```
$ go run consumer.go
2020/09/26 22:04:02 start consumer...
```

open a new terminal, send a message to rabbit mq.
```
$ go run producer.go
2020/09/26 22:22:54 send message: this is message
```
switch to the consumer and pressed `control + c`, you will see like this:
```
2020/09/26 22:04:02 start consumer...
2020/09/26 22:22:54 i: 0
2020/09/26 22:22:55 i: 1
2020/09/26 22:22:56 i: 2
2020/09/26 22:22:57 i: 3
^C2020/09/26 22:22:57 receive a shutdown signal
2020/09/26 22:22:58 i: 4
2020/09/26 22:22:59 i: 5
2020/09/26 22:23:00 i: 6
2020/09/26 22:23:01 i: 7
2020/09/26 22:23:02 i: 8
2020/09/26 22:23:03 i: 9
2020/09/26 22:23:04 mq's data: this is message
2020/09/26 22:23:04 close the consumer.
```
I pressed `control + c` at the third second, and then the program continued to execute. After the entire task was executed, the entire program was actually closed.


## New features
Set the number of concurrent

```
server := lib.NewServer(work)
server.SetGraceful(time.Second * 20)
// you can Set the number of concurrent by call this function.
server.SetConcurrency(3)
server.Run(consume)
```

