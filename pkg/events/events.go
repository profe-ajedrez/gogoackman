package events

// Listener defines a callback function to be invoked when the event is triggered
type Listener func(event Event) (bool, error)

// EventInterface defines the functions an event will have
type EventInterface interface {
	// AddEventListener Attachs a listener to an event
	AddEventListener(eventName string, callback Listener) int
	// RemoveEventListener removes a listener from an event
	RemoveEventListener(eventName string, listenerIndex int)
	// Trigger invokes all the listeners callbacks binded to an event
	Trigger(eventName string)
}

// Event is a situation wich is triggered and gives course to a callback
type Event struct {
	Name      string                `json:"name"`
	Listeners map[string][]Listener `json:"listeners"`
}

// AddEventListener Attachs a listener to an event
func (e *Event) AddEventListener(eventName string, callback Listener) int {
	e.Listeners[eventName] = append(e.Listeners[eventName], callback)
	return len(e.Listeners) - 1
}

// RemoveEventListener removes a listener from an event
func (e *Event) RemoveEventListener(eventName string, listenerIndex int) {
	if e.Listeners == nil || len(e.Listeners[eventName]) == 0 {
		panic("No listeners defined in event " + e.Name)
	}
	if listenerIndex < 0 || listenerIndex > len(e.Listeners[eventName])-1 {
		panic("Index out of bounds for listener array in " + e.Name + " event")
	}

	l := len(e.Listeners[eventName]) - 1

	e.Listeners[eventName][listenerIndex] = e.Listeners[eventName][l]
	e.Listeners[eventName] = e.Listeners[eventName][:l]
}

// Trigger invokes all the listeners callbacks binded to an event
func (e *Event) Trigger(eventName string) {
	for _, listener := range e.Listeners[eventName] {
		ok, err := listener(*e)
		if err != nil {
			panic(err)
		}
		if !ok {
			break
		}
	}
}
