package session

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type OTRSuite struct{}

var _ = Suite(&OTRSuite{})

func (s *OTRSuite) Test_GetAndWipeSymmetricKeyFor_failsIfConversationNotFound(c *C) {
	sess := &session{}
	sess.convManager = otrclient.NewConversationManager(nil, nil, "", nil, nil)

	res := sess.GetAndWipeSymmetricKeyFor(jid.Parse("some@one.org/foo"))
	c.Assert(res, IsNil)
}

type mockCommandManager struct{}

func (*mockCommandManager) ExecuteCmd(c interface{}) {}

type collectingSender struct {
	sentPeer  []jid.Any
	sentMsg   []string
	sentOTR   []bool
	returnErr error
}

func (cs *collectingSender) Send(peer jid.Any, msg string, otr bool) error {
	cs.sentPeer = append(cs.sentPeer, peer)
	cs.sentMsg = append(cs.sentMsg, msg)
	cs.sentOTR = append(cs.sentOTR, otr)

	return cs.returnErr
}

func initializeFullConversationAliceAndBob() (alicesess, bobsess *session, aliceconv, bobconv otrclient.Conversation, hook *test.Hook, alicecs, bobcs *collectingSender) {
	l, h := test.NewNullLogger()

	alicecs = &collectingSender{}
	sess := &session{
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{},
		cmdManager:    &mockCommandManager{},
		log:           l,
	}
	sess.convManager = otrclient.NewConversationManager(sess.newConversation, alicecs, "alice@one.org", sess.onOtrEventHandlerCreate, l)
	sess.privateKeys = []otr3.PrivateKey{alicePrivateKey}

	bobcs = &collectingSender{}
	bobsess = &session{
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{},
		cmdManager:    &mockCommandManager{},
		log:           l,
	}
	bobsess.convManager = otrclient.NewConversationManager(bobsess.newConversation, bobcs, "bob@one.org", bobsess.onOtrEventHandlerCreate, l)
	bobsess.privateKeys = []otr3.PrivateKey{bobPrivateKey}

	conv, _ := sess.convManager.EnsureConversationWith(jid.Parse("bob@one.org/foo"), nil)
	bobconv, _ = bobsess.convManager.EnsureConversationWith(jid.Parse("alice@one.org/bar"), nil)

	_ = conv.StartEncryptedChat()

	_, _ = bobconv.Receive([]byte(alicecs.sentMsg[0]))
	_, _ = conv.Receive([]byte(bobcs.sentMsg[0]))
	_, _ = bobconv.Receive([]byte(alicecs.sentMsg[1]))
	_, _ = conv.Receive([]byte(bobcs.sentMsg[1]))
	_, _ = bobconv.Receive([]byte(alicecs.sentMsg[2]))

	return sess, bobsess, conv, bobconv, h, alicecs, bobcs
}

func (s *OTRSuite) Test_GetAndWipeSymmetricKeyFor_returnsASymmetricKey(c *C) {
	sess, _, conv, _, _, _, _ := initializeFullConversationAliceAndBob()

	conv.(otr3.ReceivedKeyHandler).ReceivedSymmetricKey(32, []byte{0xFF, 0x00, 0x00}, []byte{0xBB, 0xAA, 0xFF})
	res := sess.GetAndWipeSymmetricKeyFor(jid.Parse("bob@one.org/foo"))
	c.Assert(res, DeepEquals, []byte{0xBB, 0xAA, 0xFF})
}

func (s *OTRSuite) Test_CreateSymmetricKeyFor_works(c *C) {
	sess, _, _, _, _, _, _ := initializeFullConversationAliceAndBob()

	res := sess.CreateSymmetricKeyFor(jid.Parse("bob@one.org/foo"))
	c.Assert(res, Not(IsNil))
}

func (s *OTRSuite) Test_CreateSymmetricKeyFor_failsIfConversationNotFound(c *C) {
	sess := &session{}
	sess.convManager = otrclient.NewConversationManager(nil, nil, "", nil, nil)

	res := sess.CreateSymmetricKeyFor(jid.Parse("some@one.org/foo"))
	c.Assert(res, IsNil)
}

func (s *OTRSuite) Test_CreateSymmetricKeyFor_failsIfNoOTRConversationIsSetUp(c *C) {
	sess, _, _, _, _, _, _ := initializeFullConversationAliceAndBob()
	_, _ = sess.convManager.EnsureConversationWith(jid.Parse("some@one.org"), nil)

	res := sess.CreateSymmetricKeyFor(jid.Parse("some@one.org"))
	c.Assert(res, IsNil)
}

func (s *OTRSuite) Test_StartSMP_works(c *C) {
	sess, _, _, _, hook, _, _ := initializeFullConversationAliceAndBob()

	sess.StartSMP(jid.R("bob@one.org/foo"), "something", "else")

	c.Assert(hook.Entries, HasLen, 0)
}

func (s *OTRSuite) Test_StartSMP_failsWithoutOTRConversation(c *C) {
	sess, _, _, _, hook, _, _ := initializeFullConversationAliceAndBob()
	_, _ = sess.convManager.EnsureConversationWith(jid.Parse("some@one.org/foo"), nil)

	sess.StartSMP(jid.R("some@one.org/foo"), "something", "else")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "cannot start SMP")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "otr: can't authenticate a peer without a secure conversation established")
}

func (s *OTRSuite) Test_StartSMP_failsWithoutConversation(c *C) {
	sess, _, _, _, hook, _, _ := initializeFullConversationAliceAndBob()

	sess.StartSMP(jid.R("some@one.org/foo"), "something", "else")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "tried to start SMP when a conversation does not exist")
}

func (s *OTRSuite) Test_FinishSMP_works(c *C) {
	sess, bobsess, conv, _, hook, _, bobcs := initializeFullConversationAliceAndBob()

	bobsess.StartSMP(jid.R("alice@one.org/bar"), "something", "else")
	_, _ = conv.Receive([]byte(bobcs.sentMsg[2]))

	sess.FinishSMP(jid.R("bob@one.org/foo"), "else")

	c.Assert(hook.Entries, HasLen, 0)
}

func (s *OTRSuite) Test_FinishSMP_failsWhenSomethingFails(c *C) {
	sess, bobsess, conv, _, hook, alicecs, bobcs := initializeFullConversationAliceAndBob()

	bobsess.StartSMP(jid.R("alice@one.org/bar"), "something", "else")
	_, _ = conv.Receive([]byte(bobcs.sentMsg[2]))

	alicecs.returnErr = errors.New("marker")
	sess.FinishSMP(jid.R("bob@one.org/foo"), "else")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "cannot provide an authentication secret for SMP")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "marker")
}

func (s *OTRSuite) Test_AbortSMP_works(c *C) {
	sess, bobsess, conv, _, hook, _, bobcs := initializeFullConversationAliceAndBob()

	bobsess.StartSMP(jid.R("alice@one.org/bar"), "something", "else")
	_, _ = conv.Receive([]byte(bobcs.sentMsg[2]))

	sess.AbortSMP(jid.R("bob@one.org/foo"))

	c.Assert(hook.Entries, HasLen, 0)
}

func (s *OTRSuite) Test_AbortSMP_failsWhenSomethingFails(c *C) {
	sess, bobsess, conv, _, hook, alicecs, bobcs := initializeFullConversationAliceAndBob()

	bobsess.StartSMP(jid.R("alice@one.org/bar"), "something", "else")
	_, _ = conv.Receive([]byte(bobcs.sentMsg[2]))

	alicecs.returnErr = errors.New("marker")
	sess.AbortSMP(jid.R("bob@one.org/foo"))

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "cannot abort SMP")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "marker")
}

func (s *OTRSuite) Test_ManuallyEndEncryptedChat_works(c *C) {
	sess, _, _, _, _, _, _ := initializeFullConversationAliceAndBob()

	e := sess.ManuallyEndEncryptedChat(jid.R("bob@one.org/foo"))
	c.Assert(e, IsNil)
}

func (s *OTRSuite) Test_ManuallyEndEncryptedChat_failsOnUnknownConversation(c *C) {
	sess, _, _, _, _, _, _ := initializeFullConversationAliceAndBob()

	e := sess.ManuallyEndEncryptedChat(jid.R("bla@one.org/foo"))
	c.Assert(e, ErrorMatches, "couldn't find conversation with.*")
}

func (s *OTRSuite) Test_terminateConversations_works(c *C) {
	sess, _, _, _, _, _, _ := initializeFullConversationAliceAndBob()

	sess.terminateConversations()
}

func (s *OTRSuite) Test_newOTRKeys_sendsNotification(c *C) {
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

	sess.newOTRKeys(jid.R("some@one.org/foo"), nil)

	<-done
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.Peer{
			Type: events.OTRNewKeys,
			From: jid.R("some@one.org/foo"),
		}})
}

func (s *OTRSuite) Test_renewedOTRKeys_sendsNotification(c *C) {
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

	sess.renewedOTRKeys(jid.R("some@one.org/foo"), nil)

	<-done
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.Peer{
			Type: events.OTRRenewedKeys,
			From: jid.R("some@one.org/foo"),
		}})
}

func (s *OTRSuite) Test_otrEnded_sendsNotification(c *C) {
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

	sess.otrEnded(jid.R("some@one.org/foo"))

	<-done
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.Peer{
			Type: events.OTREnded,
			From: jid.R("some@one.org/foo"),
		}})
}

func (s *OTRSuite) Test_listenToOtrDelayedMessageDelivery_works(c *C) {
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

	cc := make(chan int, 2)

	go func() {
		sess.listenToOtrDelayedMessageDelivery(cc, jid.R("some@one.org/foo"))
	}()

	cc <- 41
	cc <- 25
	close(cc)

	<-done
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.DelayedMessageSent{
			Peer:   jid.R("some@one.org/foo"),
			Tracer: 41,
		},
		events.DelayedMessageSent{
			Peer:   jid.R("some@one.org/foo"),
			Tracer: 25,
		},
	})
}

func (s *OTRSuite) Test_listenToOtrNotifications_works(c *C) {
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

	cc := make(chan string, 2)

	go func() {
		sess.listenToOtrNotifications(cc, jid.R("some@one.org/foo"))
	}()

	cc <- "foo"
	<-done
	cc <- "bar"
	close(cc)

	<-done
	close(ch)
	wg.Wait()

	c.Assert(nots, DeepEquals, []interface{}{
		events.Notification{
			Peer:         jid.R("some@one.org/foo"),
			Notification: "foo",
		},
		events.Notification{
			Peer:         jid.R("some@one.org/foo"),
			Notification: "bar",
		},
	})
}
