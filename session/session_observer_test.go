package session

import (
	"fmt"
	"reflect"
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

	evs := []interface{}{nil, nil, nil, nil}
	go func() {
		evs[0] = <-evc
		evs[1] = <-evc
		evs[2] = <-evc
		evs[3] = <-evc
		for _, sub := range sess.subscribers.subs {
			close(sub)
		}
		sess.subscribers.subs = sess.subscribers.subs[0:0]
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

	<-done

	c.Assert(sess.r.ToSlice(), HasLen, 0)
	c.Assert(listContains(evs, "hello"), Equals, true)
	c.Assert(listContains(evs, events.Event{Type: events.Ping}), Equals, true)
	c.Assert(listContains(evs, events.Event{Type: events.Disconnected}), Equals, true)
	c.Assert(listContains(evs, events.Event{Type: events.RosterReceived}), Equals, true)
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
			fmt.Printf("reached timeout waiting for subscribers\n")
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

	evs := []interface{}{nil, nil, nil, nil}
	go func() {
		evs[0] = <-evc
		evs[1] = <-evc
		evs[2] = <-evc
		evs[3] = <-evc
		for _, sub := range sess.subscribers.subs {
			close(sub)
		}
		sess.subscribers.subs = sess.subscribers.subs[0:0]
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

	<-done

	c.Assert(sess.r.ToSlice(), HasLen, 0)

	c.Assert(listContains(evs, "hello"), Equals, true)
	c.Assert(listContains(evs, events.Event{Type: events.Ping}), Equals, true)
	c.Assert(listContains(evs, events.Event{Type: events.ConnectionLost}), Equals, true)
	c.Assert(listContains(evs, events.Event{Type: events.RosterReceived}), Equals, true)
}

func listContains(ll []interface{}, val interface{}) bool {
	for _, v := range ll {
		if reflect.DeepEqual(v, val) {
			return true
		}
	}
	return false
}
