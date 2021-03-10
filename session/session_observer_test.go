package session

import (
	"time"

	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
	. "gopkg.in/check.v1"
)

type SessionObserverSuite struct{}

var _ = Suite(&SessionObserverSuite{})

func (s *SessionObserverSuite) Test_observe_onDisconnected(c *C) {
	done := make(chan bool)
	evc := make(chan interface{})
	evsDone := make(chan bool)

	sess := &session{
		r:                 roster.New(),
		eventsReachedZero: evsDone,
	}

	sess.Subscribe(evc)

	evs := []interface{}{nil, nil, nil}
	go func() {
		evs[0] = <-evc
		evs[1] = <-evc
		evs[2] = <-evc
	}()

	go func() {
		observe(sess)
		done <- true
	}()

	waitUntilHasSubscribers(sess, 2)

	sess.publishEvent("hello")
	<-evsDone
	sess.publish(events.Ping)
	<-evsDone

	_ = sess.r.PeerPresenceUpdate(jid.R("hello@bar.com/foo"), "available", "somewhere", "")

	c.Assert(sess.r.ToSlice(), HasLen, 1)

	sess.publish(events.Disconnected)
	<-evsDone
	close(sess.subscribers.subs[0])
	close(sess.subscribers.subs[1])

	<-done

	c.Assert(sess.r.ToSlice(), HasLen, 0)
	c.Assert(evs[0:3], DeepEquals, []interface{}{"hello", events.Event{Type: events.Ping}, events.Event{Type: events.Disconnected}})
}

func waitUntilHasSubscribers(sess *session, num int) {
	timeout := time.After(3 * time.Second)
	for {
		select {
		case <-time.After(10 * time.Millisecond):
			if sess.subscribers.subs != nil && len(sess.subscribers.subs) >= num {
				return
			}
		case <-timeout:
			return
		}
	}
}

func (s *SessionObserverSuite) Test_observe_onConnectionLost(c *C) {
	done := make(chan bool)
	evc := make(chan interface{})
	evsDone := make(chan bool)

	sess := &session{
		r:                 roster.New(),
		eventsReachedZero: evsDone,
	}

	sess.Subscribe(evc)

	evs := []interface{}{nil, nil, nil}
	go func() {
		evs[0] = <-evc
		evs[1] = <-evc
		evs[2] = <-evc
	}()

	go func() {
		observe(sess)
		done <- true
	}()

	waitUntilHasSubscribers(sess, 2)

	sess.publishEvent("hello")
	<-evsDone
	sess.publish(events.Ping)
	<-evsDone

	_ = sess.r.PeerPresenceUpdate(jid.R("hello@bar.com/foo"), "available", "somewhere", "")

	c.Assert(sess.r.ToSlice(), HasLen, 1)

	sess.publish(events.ConnectionLost)
	<-evsDone
	close(sess.subscribers.subs[0])
	close(sess.subscribers.subs[1])

	<-done

	c.Assert(sess.r.ToSlice(), HasLen, 0)
	c.Assert(evs[0:3], DeepEquals, []interface{}{"hello", events.Event{Type: events.Ping}, events.Event{Type: events.ConnectionLost}})
}
