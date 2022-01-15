package rabbit

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const EventConnect = "connect"
const EventReconnect = "reconnect"
const EventChange = "change"
const EventFailDialing = "dialing failed"
const EventFailConnect = "connection failed"
const EventFailStreamQueue = "stream queue failed"
const EventFailPushQueue = "push queue failed"
const EventRetry = "reconnection failed"
const EventListening = "listening"
const EventClosing = "closing"

const timeToResend = 5 * time.Second
const timeToReconnect = 5 * time.Second

type AckmanConsumer struct {
	pushQueue     string
	streamQueue   string
	connection    *amqp.Connection
	channel       *amqp.Channel
	Done          chan os.Signal
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
	IsConnected   bool
	Alive         bool
	Threads       int
	Wg            *sync.WaitGroup
	State         string
	events        map[string]EventCallback
}

// Connect makes an attempt to connect to RabbitMq. Returns the success.
func (a *AckmanConsumer) Connect(address string) bool {
	a.State = "CONNECTING"
	conn, err := amqp.Dial(address)
	if err != nil {
		invokeCallback(a, a.events, EventFailDialing, err)
		return false
	}

	ch, err := conn.Channel()
	if err != nil {
		invokeCallback(a, a.events, EventFailConnect, err)
		return false
	}

	ch.Confirm(false)
	_, err = ch.QueueDeclare(
		a.streamQueue,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)

	if err != nil {
		invokeCallback(a, a.events, EventFailStreamQueue, err)
		return false
	}

	_, err = ch.QueueDeclare(
		a.pushQueue,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		invokeCallback(a, a.events, EventFailPushQueue, err)
		return false
	}
	a.Change(conn, ch)
	a.IsConnected = true

	invokeCallback(a, a.events, EventConnect, nil)
	a.State = "LISTENING"
	return true
}

func (a *AckmanConsumer) Reconnect(address string) {
	for a.Alive {
		a.IsConnected = false
		fmt.Printf("Attempting to connect to rabbitMQ: %s\n", address)
		var retryCount int
		for !a.Connect(address) {
			if !a.Alive {
				return
			}
			select {
			case <-a.Done:
				return
			case <-time.After(timeToReconnect + time.Duration(retryCount)*time.Second):
				invokeCallback(a, a.events, EventRetry, nil)
				retryCount++
			}
		}

		select {
		case <-a.Done:
			invokeCallback(a, a.events, EventListening, nil)
			return
		case <-a.notifyClose:
			invokeCallback(a, a.events, EventClosing, nil)
		}
	}
}

// Change takes a new connection to the queue, and updates the channel listeners
func (a *AckmanConsumer) Change(connection *amqp.Connection, channel *amqp.Channel) {
	a.State = "CHANGING"
	a.connection = connection
	a.channel = channel
	a.notifyClose = make(chan *amqp.Error)
	a.notifyConfirm = make(chan amqp.Confirmation)
	a.channel.NotifyClose(a.notifyClose)
	a.channel.NotifyPublish(a.notifyConfirm)

	invokeCallback(a, a.events, EventChange, nil)

	a.State = "IDLE"
}

func (a *AckmanConsumer) Push(data []byte) error
func (a *AckmanConsumer) UnsafePush(data []byte) error
func (a *AckmanConsumer) Stream(cancelCtx context.Context) error
func (a *AckmanConsumer) parseEvent(msg amqp.Delivery)
func (a *AckmanConsumer) Close() error
func (a *AckmanConsumer) GetState() string
func (a *AckmanConsumer) GetSubject() interface{}

func invokeCallback(a Consumer, events map[string]EventCallback, eventName string, err error) {
	eventCallback, ok := events[eventName]
	if ok {
		eventCallback(a, eventName, err)
	}
}
