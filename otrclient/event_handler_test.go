package otrclient

import (
	"bytes"
	"errors"
	"io/ioutil"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/otr3"

	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type EventHandlerSuite struct{}

var _ = Suite(&EventHandlerSuite{})

func (s *EventHandlerSuite) Test_HandleErrorMessage_handlesAllErrorMessages(c *C) {
	handler := &EventHandler{log: log.New().WithFields(log.Fields{})}
	c.Check(string(handler.HandleErrorMessage(otr3.ErrorCodeEncryptionError)), DeepEquals, "Error occurred encrypting message.")
	c.Check(string(handler.HandleErrorMessage(otr3.ErrorCodeMessageUnreadable)), DeepEquals, "You transmitted an unreadable encrypted message.")
	c.Check(string(handler.HandleErrorMessage(otr3.ErrorCodeMessageMalformed)), DeepEquals, "You transmitted a malformed data message.")
	c.Check(string(handler.HandleErrorMessage(otr3.ErrorCodeMessageNotInPrivate)), DeepEquals, "You sent encrypted data to a peer, who wasn't expecting it.")
	c.Check(handler.HandleErrorMessage(otr3.ErrorCode(42)), IsNil)
}

func (s *EventHandlerSuite) Test_HandleSecurityEvent_HandlesAllSecurityEvents(c *C) {
	handler := &EventHandler{securityChange: -1, log: log.New().WithFields(log.Fields{})}

	handler.HandleSecurityEvent(otr3.GoneSecure)
	c.Assert(handler.securityChange, Equals, NewKeys)

	handler.HandleSecurityEvent(otr3.StillSecure)
	c.Assert(handler.securityChange, Equals, RenewedKeys)

	handler.HandleSecurityEvent(otr3.GoneInsecure)
	c.Assert(handler.securityChange, Equals, ConversationEnded)

	handler.securityChange = -1
	handler.HandleSecurityEvent(otr3.SecurityEvent(42))
	c.Assert(handler.securityChange, Equals, SecurityChange(-1))
}

func (s *EventHandlerSuite) Test_ConsumeSecurityChange_returnsTheChangeAndSetsItBack(c *C) {
	handler := &EventHandler{securityChange: RenewedKeys, log: log.New().WithFields(log.Fields{})}

	res := handler.ConsumeSecurityChange()
	c.Assert(handler.securityChange, Equals, NoChange)
	c.Assert(res, Equals, RenewedKeys)
}

func (s *EventHandlerSuite) Test_HandleSMPEvent_handlesSMPEventsAboutSecrets(c *C) {
	handler := &EventHandler{log: log.New().WithFields(log.Fields{})}
	handler.HandleSMPEvent(otr3.SMPEventAskForSecret, 72, "foo bar?")
	c.Assert(handler.securityChange, Equals, SMPSecretNeeded)
	c.Assert(handler.SmpQuestion, Equals, "foo bar?")
	c.Assert(handler.WaitingForSecret, Equals, true)

	handler = &EventHandler{log: log.New().WithFields(log.Fields{})}
	handler.HandleSMPEvent(otr3.SMPEventAskForAnswer, 72, "foo bar2?")
	c.Assert(handler.securityChange, Equals, SMPSecretNeeded)
	c.Assert(handler.SmpQuestion, Equals, "foo bar2?")
	c.Assert(handler.WaitingForSecret, Equals, true)
}

func (s *EventHandlerSuite) Test_HandleSMPEvent_handlesSMPEventsAboutSuccess(c *C) {
	handler := &EventHandler{log: log.New().WithFields(log.Fields{})}
	handler.HandleSMPEvent(otr3.SMPEventSuccess, 72, "")
	c.Assert(handler.securityChange, Equals, NoChange)

	handler.HandleSMPEvent(otr3.SMPEventSuccess, 100, "")
	c.Assert(handler.securityChange, Equals, SMPComplete)
}

func (s *EventHandlerSuite) Test_HandleSMPEvent_handlesSMPEventsAboutFailure(c *C) {
	handler := &EventHandler{log: log.New().WithFields(log.Fields{})}
	handler.HandleSMPEvent(otr3.SMPEventAbort, 72, "")
	c.Assert(handler.securityChange, Equals, SMPFailed)

	handler = &EventHandler{log: log.New().WithFields(log.Fields{})}
	handler.HandleSMPEvent(otr3.SMPEventFailure, 72, "")
	c.Assert(handler.securityChange, Equals, SMPFailed)

	handler = &EventHandler{log: log.New().WithFields(log.Fields{})}
	handler.HandleSMPEvent(otr3.SMPEventCheated, 72, "")
	c.Assert(handler.securityChange, Equals, SMPFailed)

	handler = &EventHandler{log: log.New().WithFields(log.Fields{})}
	handler.HandleSMPEvent(otr3.SMPEvent(42), 72, "")
	c.Assert(handler.securityChange, Equals, NoChange)
}

func captureLog(ll *log.Logger, f func()) string {
	ll.SetLevel(log.DebugLevel)
	buf := new(bytes.Buffer)
	// prevFlags := log.Flags()
	// log.SetFlags(0)
	ll.SetOutput(buf)
	f()
	// log.SetFlags(prevFlags)
	ll.SetOutput(ioutil.Discard)
	return buf.String()
}

func (s *EventHandlerSuite) Test_HandleMessageEvent_logsHeartbeatEvents(c *C) {
	ll := log.New()
	handler := &EventHandler{account: "me1@foo.bar", peer: jid.NR("them1@somewhere.com"), log: ll.WithField("account", "me1@foo.bar")}
	l := captureLog(ll, func() {
		handler.HandleMessageEvent(otr3.MessageEventLogHeartbeatReceived, nil, nil)
	})

	c.Assert(handler.securityChange, Equals, NoChange)
	c.Assert(l, Matches, ".*?account=me1@foo\\.bar.*?\n")
	c.Assert(l, Matches, ".*?from=them1@somewhere\\.com.*?\n")
	c.Assert(l, Matches, ".*?Heartbeat received.*?\n")

	l2 := captureLog(ll, func() {
		handler.HandleMessageEvent(otr3.MessageEventLogHeartbeatSent, nil, nil)
	})

	c.Assert(handler.securityChange, Equals, NoChange)
	c.Assert(l2, Matches, ".*?account=me1@foo\\.bar.*?\n")
	c.Assert(l2, Matches, ".*?to=them1@somewhere\\.com.*?\n")
	c.Assert(l2, Matches, ".*?Heartbeat sent.*?\n")
}

func (s *EventHandlerSuite) Test_HandleMessageEvent_logsUnrecognizedMessage(c *C) {
	ll := log.New()
	handler := &EventHandler{account: "me1@foo.bar", peer: jid.NR("them1@somewhere.com"), log: ll.WithField("account", "me1@foo.bar")}
	l := captureLog(ll, func() {
		handler.HandleMessageEvent(otr3.MessageEventReceivedMessageUnrecognized, nil, nil)
	})

	c.Assert(l, Matches, ".*?account=me1@foo\\.bar.*?\n")
	c.Assert(l, Matches, ".*?from=them1@somewhere\\.com.*?\n")
	c.Assert(l, Matches, ".*?Unrecognized OTR message received.*?\n")
}

func (s *EventHandlerSuite) Test_HandleMessageEvent_logsUnhandledEvent(c *C) {
	ll := log.New()
	handler := &EventHandler{account: "me1@foo.bar", peer: jid.NR("them1@somewhere.com"), log: ll.WithField("account", "me1@foo.bar")}
	l := captureLog(ll, func() {
		handler.HandleMessageEvent(otr3.MessageEvent(44422), nil, nil)
	})

	c.Assert(l, Matches, ".*?account=me1@foo\\.bar.*?\n")
	c.Assert(l, Matches, ".*?event=\"MESSAGE EVENT: \\(THIS SHOULD NEVER HAPPEN\\)\".*?\n")
	c.Assert(l, Matches, ".*?msg=\"Unhandled OTR3 Message Event\".*?\n")
}

func (s *EventHandlerSuite) Test_HandleMessageEvent_ignoresMessageForOtherInstance(c *C) {
	ll := log.New()
	handler := &EventHandler{account: "me1@foo.bar", peer: jid.NR("them1@somewhere.com"), log: ll}
	l := captureLog(ll, func() {
		handler.HandleMessageEvent(otr3.MessageEventReceivedMessageForOtherInstance, nil, nil)
	})

	c.Assert(l, Equals, "")
}

func (s *EventHandlerSuite) Test_HandleMessageEvent_notifiesOnSeveralMessageEvents(c *C) {
	nn := make(chan string, 1)
	defer func() {
		close(nn)
	}()

	handler := &EventHandler{account: "me2@foo.bar", peer: jid.NR("them2@somewhere.com"), notifications: nn, delays: make(map[int]bool)}
	handler.HandleMessageEvent(otr3.MessageEventEncryptionRequired, nil, nil, 123)
	c.Assert(<-nn, Equals, "Attempting to start a private conversation...")

	handler.HandleMessageEvent(otr3.MessageEventEncryptionError, nil, nil)
	c.Assert(<-nn, Equals, "An error occurred when encrypting your message. The message was not sent.")

	handler.HandleMessageEvent(otr3.MessageEventConnectionEnded, nil, nil)
	c.Assert(<-nn, Equals, "Your message was not sent, since the other person has already closed their private connection to you.")

	handler.HandleMessageEvent(otr3.MessageEventMessageReflected, nil, nil)
	c.Assert(<-nn, Equals, "We are receiving our own OTR messages. You are either trying to talk to yourself, or someone is reflecting your messages back at you.")

	handler.HandleMessageEvent(otr3.MessageEventMessageResent, nil, nil)
	c.Assert(<-nn, Equals, "The last message to the other person was resent, since we couldn't deliver the message previously.")

	handler.HandleMessageEvent(otr3.MessageEventReceivedMessageUnreadable, nil, nil)
	c.Assert(<-nn, Equals, "We received an unreadable encrypted message. It has probably been tampered with, or was sent from an older client.")

	handler.HandleMessageEvent(otr3.MessageEventReceivedMessageMalformed, nil, nil)
	c.Assert(<-nn, Equals, "We received a malformed data message.")

	handler.HandleMessageEvent(otr3.MessageEventReceivedMessageGeneralError, []byte("hmm this is weird"), nil)
	c.Assert(<-nn, Equals, "We received this error from the other person: hmm this is weird.")

	handler.HandleMessageEvent(otr3.MessageEventReceivedMessageNotInPrivate, nil, nil)
	c.Assert(<-nn, Equals, "We received an encrypted message which can't be read, since private communication is not currently turned on. You should ask your peer to repeat what they said.")

	handler.HandleMessageEvent(otr3.MessageEventReceivedMessageUnencrypted, nil, nil)
	c.Assert(<-nn, Equals, "We received a message that was transferred without encryption")
}

func (s *EventHandlerSuite) Test_HandleMessageEvent_handlesMessageEventSetupCorrectly(c *C) {
	ll := log.New()
	nn := make(chan string, 1)
	defer func() {
		close(nn)
	}()

	handler := &EventHandler{account: "me2@foo.bar", peer: jid.NR("them2@somewhere.com"), notifications: nn, log: ll.WithField("account", "me2@foo.bar")}
	l := captureLog(ll, func() {
		handler.HandleMessageEvent(otr3.MessageEventSetupError, nil, nil)
	})
	c.Assert(<-nn, Equals, "Error setting up private conversation.")
	c.Assert(l, Equals, "")

	l = captureLog(ll, func() {
		handler.HandleMessageEvent(otr3.MessageEventSetupError, nil, errors.New("hmm bla bla"))
	})
	c.Assert(<-nn, Equals, "Error setting up private conversation.")
	c.Assert(l, Matches, ".*?account=me2@foo\\.bar.*?\n")
	c.Assert(l, Matches, ".*?msg=\"Error setting up private conversation\".*?\n")
	c.Assert(l, Matches, ".*?with=them2@somewhere.com.*?\n")
}

func (s *EventHandlerSuite) Test_EventHandler_ConsumeDelayedState(c *C) {
	ev := &EventHandler{
		delays: map[int]bool{},
	}

	ev.delays[42] = true
	res := ev.ConsumeDelayedState(42)
	c.Assert(res, Equals, true)
}

func (s *EventHandlerSuite) Test_EventHandler_HandleMessageEvent_MessageEventMessageSent_addsDelay(c *C) {
	ev := &EventHandler{
		delays: map[int]bool{},
	}

	ch := make(chan int)
	ev.delayedMessageSent = ch

	wg := sync.WaitGroup{}
	wg.Add(1)

	receivedValue := 0
	go func() {
		receivedValue = <-ch
		wg.Done()
	}()

	ev.HandleMessageEvent(otr3.MessageEventMessageSent, nil, nil)
	ev.HandleMessageEvent(otr3.MessageEventMessageSent, nil, nil, 55)

	wg.Wait()

	c.Assert(receivedValue, Equals, 55)
}
