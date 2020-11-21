package lib

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type server struct {
	wg              *sync.WaitGroup
	wait            chan bool
	handler         func(delivery amqp.Delivery)
	close           bool
	graceful        bool
	gracefulTimeout time.Duration
}

func NewServer(handler func(delivery amqp.Delivery)) *server {
	return &server{
		wg:       &sync.WaitGroup{},
		wait:     make(chan bool, 1),
		handler:  handler,
		close:    false,
		graceful: false,
	}
}

func (srv *server) SetGraceful(duration time.Duration) {
	srv.graceful = true
	srv.gracefulTimeout = duration
}

func (srv *server) Run(deliveries <-chan amqp.Delivery) {
	go func() {
		for d := range deliveries {
			go srv.startHandler(d)
			// after close, will not process new message
			if srv.close {
				break
			}
		}
	}()

	if srv.graceful {
		srv.gracefulShutdown()
		//Process the rest message after it received shutdown signal
		s.wg.Wait()
	} else {
		forever := make(chan bool)
		<-forever
	}
}

func (srv *server) waitGroup() {
	srv.wg.Wait()
	srv.wait <- true
}

func (srv *server) shutdown(ctx context.Context) error {
	go srv.waitGroup()
	select {
	case <-srv.wait:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (srv *server) startHandler(delivery amqp.Delivery) {
	srv.wg.Add(1)
	defer srv.wg.Done()
	srv.handler(delivery)
}

func (srv *server) gracefulShutdown() {
	//Block until a shutdown signal is received
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)
	_ = <-ch
	log.Println("receive a shutdown signal")

	srv.close = true

	//set the max time to close
	cxt, cancel := context.WithTimeout(context.Background(), srv.gracefulTimeout)
	defer cancel()

	// graceful shutdown
	if err := srv.shutdown(cxt); err != nil {
		log.Fatalln(err)
	}
	log.Println("close the consumer.")
}
