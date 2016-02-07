package session

import "github.com/twstrike/coyim/session/events"

func observe(s *session) {
	observer := make(chan interface{})
	s.Subscribe(observer)

	for ev := range observer {
		switch t := ev.(type) {
		case events.Event:
			switch t.Type {
			case events.Disconnected, events.ConnectionLost:
				s.r.Clear()
				s.rosterReceived()
			}
		}
	}
}
