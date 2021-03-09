package session

import (
	"sync"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	. "gopkg.in/check.v1"
)

type MUCErrorsSuite struct{}

var _ = Suite(&MUCErrorsSuite{})

func (s *MUCErrorsSuite) Test_getEventErrorTypeBasedOnStanzaError_works(c *C) {
	ee := &data.StanzaError{}
	c.Assert(getEventErrorTypeBasedOnStanzaError(ee), Equals, events.MUCNoError)

	ee = &data.StanzaError{MUCNotAuthorized: &data.MUCNotAuthorized{}}
	c.Assert(getEventErrorTypeBasedOnStanzaError(ee), Equals, events.MUCNotAuthorized)

	ee = &data.StanzaError{MUCForbidden: &data.MUCForbidden{}}
	c.Assert(getEventErrorTypeBasedOnStanzaError(ee), Equals, events.MUCForbidden)

	ee = &data.StanzaError{MUCItemNotFound: &data.MUCItemNotFound{}}
	c.Assert(getEventErrorTypeBasedOnStanzaError(ee), Equals, events.MUCItemNotFound)

	ee = &data.StanzaError{MUCNotAllowed: &data.MUCNotAllowed{}}
	c.Assert(getEventErrorTypeBasedOnStanzaError(ee), Equals, events.MUCNotAllowed)

	ee = &data.StanzaError{MUCNotAcceptable: &data.MUCNotAcceptable{}}
	c.Assert(getEventErrorTypeBasedOnStanzaError(ee), Equals, events.MUCNotAcceptable)

	ee = &data.StanzaError{MUCRegistrationRequired: &data.MUCRegistrationRequired{}}
	c.Assert(getEventErrorTypeBasedOnStanzaError(ee), Equals, events.MUCRegistrationRequired)

	ee = &data.StanzaError{MUCConflict: &data.MUCConflict{}}
	c.Assert(getEventErrorTypeBasedOnStanzaError(ee), Equals, events.MUCConflict)

	ee = &data.StanzaError{MUCServiceUnavailable: &data.MUCServiceUnavailable{}}
	c.Assert(getEventErrorTypeBasedOnStanzaError(ee), Equals, events.MUCServiceUnavailable)
}

func (s *MUCErrorsSuite) Test_isMUCError_works(c *C) {
	ee := &data.StanzaError{}
	c.Assert(isMUCError(ee), Equals, false)

	ee = &data.StanzaError{MUCNotAuthorized: &data.MUCNotAuthorized{}}
	c.Assert(isMUCError(ee), Equals, true)
}

func (s *MUCErrorsSuite) Test_getEventErrorTypeBasedOnMessageError_works(c *C) {
	ee := &data.StanzaError{}
	c.Assert(getEventErrorTypeBasedOnMessageError(ee), Equals, events.MUCNoError)

	ee = &data.StanzaError{MUCNotAuthorized: &data.MUCNotAuthorized{}}
	c.Assert(getEventErrorTypeBasedOnMessageError(ee), Equals, events.MUCNoError)

	ee = &data.StanzaError{MUCForbidden: &data.MUCForbidden{}}
	c.Assert(getEventErrorTypeBasedOnMessageError(ee), Equals, events.MUCMessageForbidden)

	ee = &data.StanzaError{MUCNotAcceptable: &data.MUCNotAcceptable{}}
	c.Assert(getEventErrorTypeBasedOnMessageError(ee), Equals, events.MUCMessageNotAcceptable)
}

func (s *MUCErrorsSuite) Test_isMUCErrorPresence_works(c *C) {
	c.Assert(isMUCErrorPresence(nil), Equals, false)

	ee := &data.StanzaError{}
	c.Assert(isMUCErrorPresence(ee), Equals, false)

	ee = &data.StanzaError{MUCNotAuthorized: &data.MUCNotAuthorized{}}
	c.Assert(isMUCErrorPresence(ee), Equals, true)

	ee = &data.StanzaError{MUCForbidden: &data.MUCForbidden{}}
	c.Assert(isMUCErrorPresence(ee), Equals, true)

	ee = &data.StanzaError{MUCItemNotFound: &data.MUCItemNotFound{}}
	c.Assert(isMUCErrorPresence(ee), Equals, true)

	ee = &data.StanzaError{MUCNotAllowed: &data.MUCNotAllowed{}}
	c.Assert(isMUCErrorPresence(ee), Equals, true)

	ee = &data.StanzaError{MUCNotAcceptable: &data.MUCNotAcceptable{}}
	c.Assert(isMUCErrorPresence(ee), Equals, true)

	ee = &data.StanzaError{MUCRegistrationRequired: &data.MUCRegistrationRequired{}}
	c.Assert(isMUCErrorPresence(ee), Equals, true)

	ee = &data.StanzaError{MUCConflict: &data.MUCConflict{}}
	c.Assert(isMUCErrorPresence(ee), Equals, true)

	ee = &data.StanzaError{MUCServiceUnavailable: &data.MUCServiceUnavailable{}}
	c.Assert(isMUCErrorPresence(ee), Equals, true)
}

func (s *MUCErrorsSuite) Test_mucManager_publishMUCError_works(c *C) {
	sess := &session{}
	ch := make(chan interface{})
	waiting := make(chan bool)
	sess.eventsReachedZero = waiting
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

	m := &mucManager{
		publishEvent: sess.publishEvent,
	}

	m.publishMUCError(jid.ParseFull("hello@goodbye.com/foo"), &data.StanzaError{MUCForbidden: &data.MUCForbidden{}})
	<-waiting

	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.MUCError{Nickname: "foo", Room: jid.ParseBare("hello@goodbye.com"), ErrorType: events.MUCForbidden},
	})
}
func (s *MUCErrorsSuite) Test_mucManager_publishMUCMessageError_works(c *C) {
	sess := &session{}
	ch := make(chan interface{})
	waiting := make(chan bool)
	sess.eventsReachedZero = waiting
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

	m := &mucManager{
		publishEvent: sess.publishEvent,
	}

	m.publishMUCMessageError(jid.ParseBare("hello@goodbye.com"), &data.StanzaError{MUCForbidden: &data.MUCForbidden{}})
	<-waiting

	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.MUCError{Nickname: "", Room: jid.ParseBare("hello@goodbye.com"), ErrorType: events.MUCMessageForbidden},
	})
}
