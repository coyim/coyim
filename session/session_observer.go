package session

func observe(s *Session) {
	observer := make(chan interface{})
	s.Subscribe(observer)

	for ev := range observer {
		switch t := ev.(type) {
		case Event:
			switch t.Type {
			case Disconnected, ConnectionLost:
				s.R.Clear()
				s.rosterReceived()
			}
		}
	}
}
