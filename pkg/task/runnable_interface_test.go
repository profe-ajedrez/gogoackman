package task

import "testing"

type runnable struct{}

func (r *runnable) Do(payload []byte)   {}
func (r *runnable) Undo(payload []byte) {}

func TestRunableInterface(t *testing.T) {
	runner := runnable{}
	runner.Do([]byte{})
}
