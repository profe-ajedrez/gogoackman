package rbmq

import (
	"log"
	"testing"
	"time"
)

type runnable struct{}

func (r *runnable) Do(payload []byte) {
	log.Println("Consumiendo " + string(payload))
}

func (r *runnable) Undo(payload []byte) {
	log.Println("Desconsumiendo " + string(payload))
}

func TestGogoAcmanConsumer(t *testing.T) {
	rc := RabbiterConfig{
		Schema:         "amqp",
		Username:       "test",
		Password:       "test",
		Host:           "ec2-13-57-35-233.us-west-1.compute.amazonaws.com",
		Port:           "5672",
		VHost:          "import",
		ConnectionName: "estulapido",
	}

	rbt := NewGoGoRabbit(rc)
	if err := rbt.Connect(); err != nil {
		log.Fatalln("unable to connect to rabbit", err)
	}

	cc := ConsumerConfig{
		ExchangeName: "amq",
		ExchangeType: "direct",
		RoutingKey:   "create",
		QueueName:    "garganta_fifada",
		ConsumerName: "consuma_test",
		WorkerQty:    5,
		PreloadLevel: 5,
		Reconnect: struct {
			MaxAttempt int
			Interval   time.Duration
		}{},
	}

	callbacks := Callbacks{
		Task: &runnable{},
		OnDisconnection: func(errorStr string) {
			log.Println("On Disconnection " + errorStr)
		},
		OnFailedAttemp: func(errorMsg string, attemp int) {
			log.Println("On Failed Attemp " + errorMsg)
			log.Println(attemp)
		},
		OnGiveUp: func() {
			log.Println("On Give up ")
		},
		OnErrorAtConsumeMessage: func(workerId int, errorStr string, body []byte) {
			log.Println("On Error At Consume " + string(body) + " " + errorStr)
			log.Println(workerId)
		},
	}

	cc.Reconnect.MaxAttempt = 60
	cc.Reconnect.Interval = 3 * time.Second

	csm := NewGogoAckman(cc, rbt)

	if err := csm.Start(callbacks); err != nil {
		log.Fatalln("unable to start consumer", err)
	}
	//

	select {}
}
