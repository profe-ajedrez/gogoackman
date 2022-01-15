package rabbit

import (
	"context"

	"github.com/streadway/amqp"
)

type EventCallback func(target Consumer, event string, err error)

type Events interface {
	On(event string, callback EventCallback)
}

type Consumer interface {
	Connect(address string) bool
	Reconnect(address string)
	Change(connection *amqp.Connection, channel *amqp.Channel)
	Push(data []byte) error
	UnsafePush(data []byte) error
	Stream(cancelCtx context.Context) error
	parseEvent(msg amqp.Delivery)
	Close() error
	GetState() string
	GetSubject() interface{}
}
