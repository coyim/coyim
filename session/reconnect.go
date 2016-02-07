package session

import (
	"math/rand"
	"time"
)

func checkReconnect(s *Session) {
	for {
		<-time.After(time.Duration(rand.Int31n(7643)) * time.Millisecond)

		if s.IsDisconnected() && s.wantToBeOnline {
			s.connector.Connect()
		}
	}
}
