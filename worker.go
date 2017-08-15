package rabbit

import (
	"github.com/hippoai/goutil"
	"github.com/streadway/amqp"
)

type WorkerQueue struct {
	Rabbit       *Rabbit
	RoutingKey   string
	ExchangeName string
	MaxPerWorker int
}

func NewWorkerQueue(myRabbit *Rabbit, routingKey, exchangeName string, maxPerWorker int) *WorkerQueue {
	return &WorkerQueue{
		Rabbit:       myRabbit,
		RoutingKey:   routingKey,
		ExchangeName: exchangeName,
		MaxPerWorker: maxPerWorker,
	}
}

// Send a message to the worker queue
func (wq *WorkerQueue) Send(msg []byte) error {

	var err error
	defer wq.Rabbit.Close()

	// Declare the exchange
	err = wq.Rabbit.Channel.ExchangeDeclare(
		wq.ExchangeName, "direct", true,
		false, false, false, nil,
	)
	if err != nil {
		return err
	}

	// Publish the message
	err = wq.Rabbit.Channel.Publish(
		wq.ExchangeName, wq.RoutingKey,
		false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// Listen to the messages on this worker queue
func (wq *WorkerQueue) Listen(cb func([]byte) error, cbError func(err error)) error {

	consumerKey := goutil.UuidV4()

	var err error
	defer wq.Rabbit.Close()

	// Declare the exchange
	err = wq.Rabbit.Channel.ExchangeDeclare(
		wq.ExchangeName, "direct", true,
		false, false, false, nil,
	)
	if err != nil {
		return err
	}

	// Declare the queue with its routing key
	_, err = wq.Rabbit.Channel.QueueDeclare(
		wq.RoutingKey, true,
		false, false, false, nil,
	)
	if err != nil {
		return err
	}

	// Make sure we don't have more than maxPerWorker messages per worker
	err = wq.Rabbit.Channel.Qos(
		wq.MaxPerWorker,
		0,
		false,
	)

	// Bind the queue to the exchange, with its routing key
	err = wq.Rabbit.Channel.QueueBind(
		wq.RoutingKey, wq.RoutingKey, wq.ExchangeName,
		false, nil,
	)
	if err != nil {
		return err
	}

	// Create channel to consume messages
	msgs, err := wq.Rabbit.Channel.Consume(
		wq.RoutingKey, consumerKey, true,
		false, false, false, nil,
	)

	// Read messages forever
	for d := range msgs {
		err := cb(d.Body)
		if err != nil {
			cbError(err)
		}
	}

	return nil

}
