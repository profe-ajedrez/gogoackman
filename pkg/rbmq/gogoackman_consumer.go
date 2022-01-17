package rbmq

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// gogoAckman Our implementation of Consumer
type gogoAckman struct {
	config ConsumerConfig
	Rabbit Rabbiter
}

// NewGogoAckman returns a consumer instance.
func NewGogoAckman(config ConsumerConfig, rabbit Rabbiter) Consumer {
	return &gogoAckman{
		config: config,
		Rabbit: rabbit,
	}
}

// Start runs the workers
func (a *gogoAckman) Start(actionCallback ConsumerCallback) error {
	con, err := a.Rabbit.Connection()
	if err != nil {
		return err
	}

	go a.closedConnectionListener(con.NotifyClose(make(chan *amqp.Error)), actionCallback)

	chn, err := con.Channel()
	if err != nil {
		return err
	}

	if err := chn.ExchangeDeclare(
		a.config.ExchangeName, // name
		a.config.ExchangeType, // type
		true,                  // durable
		false,                 // autoDelete
		false,                 // internal
		false,                 // noWait
		nil,                   // args
	); err != nil {
		return err
	}

	if _, err := chn.QueueDeclare(
		a.config.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-queue-mode": "lazy"}, // See https://www.rabbitmq.com/lazy-queues.html to understand lazy queues
	); err != nil {
		return err
	}

	if err := chn.QueueBind(
		a.config.QueueName,
		a.config.RoutingKey,
		a.config.ExchangeName,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := chn.Qos(a.config.PreloadLevel, 0, false); err != nil {
		return err
	}

	for i := 1; i <= a.config.WorkerQty; i++ {
		id := i
		go a.consume(chn, id, actionCallback)
	}

	return nil
}

// closedConnectionListener attemps to reconnect to the server and tries to reopen channels for a time
// if connection is closed
func (a *gogoAckman) closedConnectionListener(closed <-chan *amqp.Error, actionCallback ConsumerCallback) {
	log.Println("INFO: Watching closed connection")

	// If you do not want to reconnect in the case of manual disconnection
	// via RabbitMQ UI or Server restart, handle `amqp.ConnectionForced`
	// error code.
	err := <-closed
	if err != nil {
		log.Println("INFO: Closed connection:", err.Error())

		var i int

		for i = 0; i < a.config.Reconnect.MaxAttempt; i++ {
			log.Println("INFO: Attempting to reconnect")

			if err := a.Rabbit.Connect(); err == nil {
				log.Println("INFO: Reconnected")

				if err := a.Start(actionCallback); err == nil {
					break
				}
			}

			time.Sleep(a.config.Reconnect.Interval)
		}

		if i == a.config.Reconnect.MaxAttempt {
			log.Println("CRITICAL: Giving up reconnecting")

			return
		}
	} else {
		log.Println("INFO: Connection closed normally, will not reconnect")
		os.Exit(0)
	}
}

func (a *gogoAckman) consume(channel *amqp.Channel, id int, actionCallback ConsumerCallback) {
	msgs, err := channel.Consume(
		a.config.QueueName,
		fmt.Sprintf("%s (%d/%d)", a.config.ConsumerName, id, a.config.WorkerQty),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(fmt.Sprintf("CRITICAL: Unable to start consumer (%d/%d)", id, a.config.WorkerQty))
		return
	}

	log.Println("[", id, "] Running ...")
	log.Println("[", id, "] Press CTRL+C to exit ...")

	for msg := range msgs {
		log.Println("[", id, "] Consumed:", string(msg.Body))

		if err := msg.Ack(false); err != nil {
			// TODO: Should DLX the message
			log.Println("unable to acknowledge the message, dropped", err)
		}
	}

	log.Println("[", id, "] Exiting ...")
}
