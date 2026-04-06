package session

import "github.com/coyim/coyim/session/events"

func observe(s *session) {
	observer := make(chan events.Is)
	s.Subscribe(observer)

	for {
		select {
		case <-s.ctx.Done():
			return
		case ev, ok := <-observer:
			if !ok {
				return
			}
			switch t := ev.(type) {
			case events.Event:
				switch t.Type {
				case events.Disconnected, events.ConnectionLost:
					s.r.Clear()
					s.rosterUpdated()
				}
			}
		}

	}
}
