package session

import (
	"sync"
	"time"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"

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
			Type: events.Disconnected,
		})
	case <-time.After(100 * time.Millisecond):
		c.Error("Did not receive expected notification")
	}
}

func (s *SessionEventSuite) Test_publish_doesNotBlockIfThereIsNoSubscriber(c *C) {
	session := &session{}
	session.publish(events.Connected)
}

func (s *SessionEventSuite) Test_session_unsubscribe_works(c *C) {
	sess := &session{}

	c1 := make(chan interface{})
	c2 := make(chan interface{})

	sess.Subscribe(c1)
	c.Assert(sess.subscribers.subs, HasLen, 1)

	sess.unsubscribe(c2)
	c.Assert(sess.subscribers.subs, HasLen, 1)

	sess.unsubscribe(c1)
	c.Assert(sess.subscribers.subs, HasLen, 0)
}

func (s *SessionEventSuite) Test_session_publishEventTo_unsubscribesOnFailure(c *C) {
	sess := &session{}

	c1 := make(chan interface{})

	sess.Subscribe(c1)

	close(c1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	sess.publishEventTo(c1, nil, &wg)

	c.Assert(sess.subscribers.subs, HasLen, 0)
}

func (s *SessionEventSuite) Test_session_PublishEvent_works(c *C) {
	ch := make(chan interface{})
	done := make(chan bool)
	sess := &session{
		eventsReachedZero: done,
	}

	sess.subscribers.subs = append(sess.subscribers.subs, ch)

	var wg sync.WaitGroup
	wg.Add(1)
	var nots []interface{}
	go func() {
		for n := range ch {
			nots = append(nots, n)
		}
		wg.Done()
	}()

	go func() {
		sess.PublishEvent("hello")
	}()

	<-done
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{"hello"})

}

func (s *SessionEventSuite) Test_session_publishSMPEvent_works(c *C) {
	ch := make(chan interface{})
	done := make(chan bool)
	sess := &session{
		eventsReachedZero: done,
	}

	sess.subscribers.subs = append(sess.subscribers.subs, ch)

	var wg sync.WaitGroup
	wg.Add(1)
	var nots []interface{}
	go func() {
		for n := range ch {
			nots = append(nots, n)
		}
		wg.Done()
	}()

	go func() {
		sess.publishSMPEvent(events.Failure, jid.R("foo@bar.com/qux"), "something")
	}()

	<-done
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.SMP{
			Type: events.Failure,
			From: jid.R("foo@bar.com/qux"),
			Body: "something",
		},
	})
}
