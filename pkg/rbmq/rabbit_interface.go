package rbmq

import "github.com/streadway/amqp"

// RabbiterConfig stores the config to the rabbitmq connection instance
type RabbiterConfig struct {
	Schema         string
	Username       string
	Password       string
	Host           string
	Port           string
	VHost          string
	ConnectionName string
}

// Rabbiter represents the rabbitmq connection instance
type Rabbiter interface {
	Connect() error
	Connection() (*amqp.Connection, error)
	Channel() (*amqp.Channel, error)
}
