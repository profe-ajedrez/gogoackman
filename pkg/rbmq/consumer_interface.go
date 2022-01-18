package rbmq

import (
	"time"

	"github.com/profe-ajedrez/gogoackman/pkg/task"
	"github.com/streadway/amqp"
)

type (
	// ConsumerCallback callback with the task to do by the consumer
	ConsumerCallback func(workerId int, body []byte)
	// OnDisconnection what to do on disconnection
	OnDisconnection func(errorStr string)
	// OnGiveUp what to do whe giving up a connection
	OnGiveUp func()
	// OnErrorAtConsumeMessage What to do on error at comsume of a message
	OnErrorAtConsumeMessage func(workerId int, errorStr string, body []byte)

	// Callbacks struct to pass the callback parameters to the consumer
	Callbacks struct {
		Task                    task.RunnableInterface
		OnDisconnection         OnDisconnection
		OnFailedAttemp          func(errorMsg string, attemp int)
		OnGiveUp                OnGiveUp
		OnErrorAtConsumeMessage OnErrorAtConsumeMessage
	}
)

// ConsumerConfig stores the config to be passed to the created consumers
type ConsumerConfig struct {
	// Exchanges are message routing agents, defined by the virtual host within RabbitMQ.
	// An exchange is responsible for routing the messages to different queues
	// with the help of header attributes, bindings, and routing keys
	//  @watch https://youtu.be/o8eU5WiO8fw
	ExchangeName string
	ExchangeType string
	// a message attribute the exchange looks at when deciding how to route the message
	// to queues (depending on exchange type) it must be a list of words, delimited by dots.
	// The words can be anything, but usually they specify some features connected to the message.
	// A few valid routing key examples: "stock.usd.nyse", "nyse.vmw", "quick.orange.rabbit"
	// @see https://www.cloudamqp.com/blog/part4-rabbitmq-for-beginners-exchanges-routing-keys-bindings.html
	RoutingKey string
	// For more insigths about a queue @see https://www.rabbitmq.com/queues.html
	QueueName    string
	ConsumerName string
	// How many workers (consumers) this connection will spawn
	WorkerQty int
	// PreloadLevel limits the number of aknowledged messages in a queue
	// see more in https://www.rabbitmq.com/consumer-prefetch.html
	PreloadLevel int
	Reconnect    struct {
		MaxAttempt int
		Interval   time.Duration
	}
}

// Consumer is an interface wich represents a connection to a RabbitMq instance to consume messages
type Consumer interface {
	Start(callbacks Callbacks) error
	closedConnectionListener(closed <-chan *amqp.Error, callbacks Callbacks)
	consume(channel *amqp.Channel, id int, callbacks Callbacks)
}
