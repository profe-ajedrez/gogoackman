package task

// RunnableInterface binds the Do func
type RunnableInterface interface {
	Do(payload []byte)
	Undo(payload []byte)
}
