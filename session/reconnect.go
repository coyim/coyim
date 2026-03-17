package session

import (
	"math/rand"
	"time"
)

var calculateRandomDelay = func() time.Duration {
	return 10*time.Second + time.Duration(rand.Int31n(7643))*time.Millisecond
}

var randomDelayChannel = func() <-chan time.Time {
	return time.After(calculateRandomDelay())
}

func checkReconnect(s *session) {
	_, cont := <-randomDelayChannel()
	for cont {
		s.wantToBeOnlineLock.Lock()
		want := s.wantToBeOnline
		s.wantToBeOnlineLock.Unlock()

		if s.IsDisconnected() && want {
			s.connector.Connect()
		}

		_, cont = <-randomDelayChannel()
	}
}
