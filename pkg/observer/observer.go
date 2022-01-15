package observer

type Subject interface {
	Register(observer Observer)
	Unregister(observer Observer)
	NotifyAll()
}

type Observer interface {
	Update(interface{})
	GetID() string
}
