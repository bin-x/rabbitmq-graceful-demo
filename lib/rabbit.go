package lib

import (
	"errors"
	"github.com/streadway/amqp"
)

type MqClient struct {
	conn *amqp.Connection
}

func (m *MqClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}
func (m *MqClient) Connect(addr string) error {
	var err error
	m.conn, err = amqp.Dial(addr)
	if err != nil {
		return errors.New("connect error: " + err.Error())
	}
	return nil
}

func (m *MqClient) Publish(body []byte, exchangeName string, exchangeType string, routingKey string) error {
	if m.conn == nil {
		return errors.New("connect error")
	}
	ch, err := m.conn.Channel() // Get a channel from the connection
	if err != nil {
		return errors.New("get channel error" + err.Error())
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		true,         // noWait
		nil,          // arguments
	)
	if err != nil {
		return errors.New("declare exchange error: " + err.Error())
	}

	err = ch.Publish( // Publishes a message onto the queue.
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Body: body, // Our JSON body as []byte
		})

	if err != nil {
		return errors.New("publish error: " + err.Error())
	}

	return nil
}

func (m *MqClient) Consume(exchangeName string, exchangeType string, consumerName string, queueName string, bindingKey string) (<-chan amqp.Delivery, error) {
	if m.conn == nil {
		return nil, errors.New("connect error")
	}
	ch, err := m.conn.Channel()
	if err != nil {
		return nil, errors.New("get channel error" + err.Error())
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		true,         // noWait
		nil,          // arguments
	)
	if err != nil {
		return nil, errors.New("declare exchange error: " + err.Error())
	}

	_, err = ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		true,      // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.New("declare queue error: " + err.Error())
	}
	err = ch.QueueBind(
		queueName,    // name of the queue
		bindingKey,   // bindingKey
		exchangeName, // sourceExchange
		true,         // noWait
		nil,          // arguments
	)
	if err != nil {
		return nil, errors.New("bind queue error: " + err.Error())
	}

	consume, err := ch.Consume(
		queueName,    // queue
		consumerName, // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		true,         // no-wait
		nil,          // args
	)

	if err != nil {
		return nil, errors.New("get consume error：" + err.Error())
	}
	return consume, nil
}
