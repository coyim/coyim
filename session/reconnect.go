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
		if s.IsDisconnected() && s.wantToBeOnline {
			s.connector.Connect()
		}

		_, cont = <-randomDelayChannel()
	}
}
