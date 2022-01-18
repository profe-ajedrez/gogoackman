package rbmq

import (
	"log"
	"testing"
)

type runnable struct{}

func (r *runnable) Do(payload []byte) {
	log.Println("Consumiendo " + string(payload))
}

func (r *runnable) Undo(payload []byte) {
	log.Println("Desconsumiendo " + string(payload))
}

func TestGogoAcmanConsumer(t *testing.T) {
}
