package rabbit

import (
	"fmt"

	"github.com/hippoai/env"
	"github.com/streadway/amqp"
)

type DefaultFunction func(args ...interface{}) (interface{}, error)

type Rabbit struct {
	Connection        *amqp.Connection
	Channel           *amqp.Channel
	OnCloseConnection map[string]DefaultFunction
	OnCloseChannel    map[string]DefaultFunction
}

func NewRabbitFromEnv() (*Rabbit, error) {

	parsed, err := env.Parse(ENV_RMQ_HOST, ENV_RMQ_PASSWORD, ENV_RMQ_PORT, ENV_RMQ_USERNAME)
	if err != nil {
		return nil, err
	}

	return NewRabbit(
		parsed[ENV_RMQ_USERNAME], parsed[ENV_RMQ_PASSWORD],
		parsed[ENV_RMQ_HOST], parsed[ENV_RMQ_PORT],
	)

}
func NewRabbit(username, password, host, port string) (*Rabbit, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s",
		username, password,
		host, port,
	))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Rabbit{
		Connection:        conn,
		Channel:           ch,
		OnCloseChannel:    map[string]DefaultFunction{},
		OnCloseConnection: map[string]DefaultFunction{},
	}, nil
}

// Close closes the connection and channel
func (rabbit *Rabbit) Close() error {
	var err error
	err = rabbit.Connection.Close()
	if err != nil {
		return err
	}

	err = rabbit.Channel.Close()
	if err != nil {
		return err
	}

	return nil
}

// SetOnCloseConnection
func (rabbit *Rabbit) SetOnCloseConnection(key string, f DefaultFunction) {
	rabbit.OnCloseConnection[key] = f
}

// SetOnCloseChannel
func (rabbit *Rabbit) SetOnCloseChannel(key string, f DefaultFunction) {
	rabbit.OnCloseChannel[key] = f
}
