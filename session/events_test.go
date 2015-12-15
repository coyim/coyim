package session

import (
	"time"

	. "gopkg.in/check.v1"
)

type SessionEventSuite struct{}

var _ = Suite(&SessionEventSuite{})

func (s *SessionEventSuite) Test_publish_notifiesWithEvents(c *C) {
	observer := make(chan interface{}, 1)

	session := &Session{}
	session.Subscribe(observer)
	session.publish(Disconnected)

	select {
	case e := <-observer:
		c.Assert(e, DeepEquals, Event{
			Type:    Disconnected,
			Session: session,
		})
	case <-time.After(1 * time.Millisecond):
		c.Error("Did not receive expected notification")
	}
}

func (s *SessionEventSuite) Test_publish_doesNotBlockIfThereIsNoSubscriber(c *C) {
	session := &Session{}
	session.publish(Connected)
}
