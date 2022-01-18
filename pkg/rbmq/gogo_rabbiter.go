package rbmq

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type gogoRabbit struct {
	config     RabbiterConfig
	connection *amqp.Connection
}

// NewGoGoRabbit returns a RabbitMQ connection instance.
func NewGoGoRabbit(config RabbiterConfig) Rabbiter {
	return &gogoRabbit{
		config: config,
	}
}

// Connect attemps to connect to a rabbitMQ server
func (g *gogoRabbit) Connect() error {
	if g.connection == nil || g.connection.IsClosed() {
		con, err := amqp.DialConfig(fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s",
			g.config.Schema,
			g.config.Username,
			g.config.Password,
			g.config.Host,
			g.config.Port,
			g.config.VHost,
		), amqp.Config{Properties: amqp.Table{"connection_name": g.config.ConnectionName}})
		if err != nil {
			return err
		}
		g.connection = con
	}

	return nil
}

// Connection returns the current connection
func (g *gogoRabbit) Connection() (*amqp.Connection, error) {
	if g.connection == nil || g.connection.IsClosed() {
		return nil, errors.New("connection is not open")
	}

	return g.connection, nil
}

// Channel returns a instance of *amqp.Channel
func (g *gogoRabbit) Channel() (*amqp.Channel, error) {
	chn, err := g.connection.Channel()
	if err != nil {
		return nil, err
	}

	return chn, nil
}
