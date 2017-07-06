package rabbit

import "github.com/streadway/amqp"

func TextQueue(ch *amqp.Channel) (*amqp.Queue, error) {

	q, err := ch.QueueDeclare(
		RMQ_QUEUE_NAME_TEXT,
		true,
		false,
		false,
		false,
		nil,
	)

	return &q, err

}
