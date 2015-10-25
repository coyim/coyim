package session

import (
	"time"

	. "gopkg.in/check.v1"
)

type SessionEventSuite struct{}

var _ = Suite(&SessionEventSuite{})

func (s *SessionEventSuite) Test_publish_notifiesWithEvents(c *C) {
	subs := make(chan Event)

	session := &Session{}
	session.Subscribe(subs)

	go session.publish(Connected)

	select {
	case e := <-subs:
		c.Assert(e, DeepEquals, Event{
			EventType: Connected,
			Session:   session,
			From:      "",
		})
	case <-time.After(1 * time.Second):
		c.Error("Did not receive expected notification")
	}

	go session.publish(Disconnected)

	select {
	case e := <-subs:
		c.Assert(e, DeepEquals, Event{
			EventType: Disconnected,
			Session:   session,
			From:      "",
		})
	case <-time.After(1 * time.Second):
		c.Error("Did not receive expected notification")
	}
}

func (s *SessionEventSuite) Test_publish_doesNotBlockIfThereIsNoSubscriber(c *C) {
	session := &Session{}
	session.publish(Connected)
}
