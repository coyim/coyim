package session

import (
	"sync"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type EncryptionSuite struct{}

var _ = Suite(&EncryptionSuite{})

func (s *EncryptionSuite) Test_session_notify_publishesANotification(c *C) {
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

	sess.notify(jid.Parse("foo1@exmaple.org/afgd"), "something happened")
	<-waiting
	sess.notify(jid.Parse("some@thing.com"), "bla")
	<-waiting
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.Notification{Peer: jid.Parse("foo1@exmaple.org/afgd"), Notification: "something happened"},
		events.Notification{Peer: jid.Parse("some@thing.com"), Notification: "bla"},
	})
}

func (s *EncryptionSuite) Test_session_processEncryption_handlesOtrWithLogging(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)
	sess := &session{
		log: l,
	}

	sess.processEncryption(jid.Parse("some@one.org/foo"), &data.Encryption{Namespace: otrEncryptionNamespace})

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.DebugLevel)
	c.Assert(hook.LastEntry().Message, Equals, "got message marked with OTR encryption tag (XEP-0380)")
	c.Assert(hook.LastEntry().Data["peer"], Equals, "some@one.org/foo")
}

func (s *EncryptionSuite) Test_session_processEncryption_handlesOtherKnownAlgorithm(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	ch := make(chan interface{})
	done := make(chan bool)
	sess := &session{
		log:               l,
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

	sess.processEncryption(jid.Parse("some@one.org/foo"), &data.Encryption{Namespace: "urn:xmpp:openpgp:0"})

	<-done
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.Notification{Peer: jid.Parse("some@one.org/foo"), Notification: "We received a message encrypted with OpenPGP for XMPP - sadly CoyIM does not support this algorithm. Please let your contact know to encrypt using OTR, nothing else, to communicate with you."},
	})

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.InfoLevel)
	c.Assert(hook.LastEntry().Message, Equals, "got message marked with unknown encryption tag (XEP-0380)")
	c.Assert(hook.LastEntry().Data["name"], Equals, "OpenPGP for XMPP")
	c.Assert(hook.LastEntry().Data["namespace"], Equals, "urn:xmpp:openpgp:0")
	c.Assert(hook.LastEntry().Data["peer"], Equals, "some@one.org/foo")
}

func (s *EncryptionSuite) Test_session_processEncryption_handlesUnknownAlgorithm(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	ch := make(chan interface{})
	done := make(chan bool)
	sess := &session{
		log:               l,
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

	sess.processEncryption(jid.Parse("some@one.org/foo"), &data.Encryption{Namespace: "urn:xmpp:something:weird", Name: "stuffy"})

	<-done
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.Notification{Peer: jid.Parse("some@one.org/foo"), Notification: "We received a message encrypted with stuffy - sadly CoyIM does not support this algorithm. Please let your contact know to encrypt using OTR, nothing else, to communicate with you."},
	})

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.InfoLevel)
	c.Assert(hook.LastEntry().Message, Equals, "got message marked with unknown encryption tag (XEP-0380)")
	c.Assert(hook.LastEntry().Data["name"], Equals, "stuffy")
	c.Assert(hook.LastEntry().Data["namespace"], Equals, "urn:xmpp:something:weird")
	c.Assert(hook.LastEntry().Data["peer"], Equals, "some@one.org/foo")
}
