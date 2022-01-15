package observer

type SimpleSubject struct {
	observers []Observer
}

func (s *SimpleSubject) Register(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *SimpleSubject) Unregister(observer Observer) {
	success := false
	for index, obs := range s.observers {
		success = obs.GetID() == observer.GetID()
		if success {
			s.observers = remove(index, s.observers)
			break
		}
	}
}

func (s *SimpleSubject) NotifyAll() {
	for _, obs := range s.observers {
		obs.Update(s)
	}
}

func remove(index int, obs []Observer) []Observer {
	l := len(obs)
	obs[index] = obs[l-1]
	return obs[:l-1]
}
