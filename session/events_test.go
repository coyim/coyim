package session

import (
	"time"

	"github.com/twstrike/coyim/session/events"

	. "gopkg.in/check.v1"
)

type SessionEventSuite struct{}

var _ = Suite(&SessionEventSuite{})

func (s *SessionEventSuite) Test_publish_notifiesWithEvents(c *C) {
	observer := make(chan interface{}, 1)

	session := &session{}
	session.Subscribe(observer)
	session.publish(events.Disconnected)

	select {
	case e := <-observer:
		c.Assert(e, DeepEquals, events.Event{
			Type:    events.Disconnected,
			Session: session,
		})
	case <-time.After(1 * time.Millisecond):
		c.Error("Did not receive expected notification")
	}
}

func (s *SessionEventSuite) Test_publish_doesNotBlockIfThereIsNoSubscriber(c *C) {
	session := &session{}
	session.publish(events.Connected)
}
