package session

import (
	"io"
	"reflect"
	"time"

	gocheck "gopkg.in/check.v1"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/otr3"
)

type mockSessionEventHandler struct {
	info                func(string)
	warn                func(string)
	alert               func(string)
	rosterReceived      func(access.Session)
	iqReceived          func(uid string)
	newOTRKeys          func(from string, conversation *otr3.Conversation)
	otrEnded            func(uid string)
	messageReceived     func(s access.Session, from string, timestamp time.Time, encrypted bool, message []byte)
	processPresence     func(from, to, show, status string, gone bool)
	subscriptionRequest func(s access.Session, uid string)
	subscribed          func(account, peer string)
	unsubscribe         func(account, peer string)
	disconnected        func()
	registerCallback    func(title, instructions string, fields []interface{}) error
}

func (m *mockSessionEventHandler) Info(v string) {
	if m.info != nil {
		m.info(v)
	}
}

func (m *mockSessionEventHandler) Warn(v string) {
	if m.warn != nil {
		m.warn(v)
	}
}

func (m *mockSessionEventHandler) Alert(v string) {
	if m.alert != nil {
		m.alert(v)
	}
}

func (m *mockSessionEventHandler) RosterReceived(s access.Session) {
	if m.rosterReceived != nil {
		m.rosterReceived(s)
	}
}

func (m *mockSessionEventHandler) IQReceived(uid string) {
	if m.iqReceived != nil {
		m.iqReceived(uid)
	}
}

func (m *mockSessionEventHandler) NewOTRKeys(from string, conversation *otr3.Conversation) {
	if m.newOTRKeys != nil {
		m.newOTRKeys(from, conversation)
	}
}

func (m *mockSessionEventHandler) OTREnded(uid string) {
	if m.otrEnded != nil {
		m.otrEnded(uid)
	}
}

func (m *mockSessionEventHandler) MessageReceived(s access.Session, from string, timestamp time.Time, encrypted bool, message []byte) {
	if m.messageReceived != nil {
		m.messageReceived(s, from, timestamp, encrypted, message)
	}
}

func (m *mockSessionEventHandler) ProcessPresence(from, to, show, status string, gone bool) {
	if m.processPresence != nil {
		m.processPresence(from, to, show, status, gone)
	}
}

func (m *mockSessionEventHandler) SubscriptionRequest(s access.Session, uid string) {
	if m.subscriptionRequest != nil {
		m.subscriptionRequest(s, uid)
	}
}

func (m *mockSessionEventHandler) Subscribed(account, peer string) {
	if m.subscribed != nil {
		m.subscribed(account, peer)
	}
}

func (m *mockSessionEventHandler) Unsubscribe(account, peer string) {
	if m.unsubscribe != nil {
		m.unsubscribe(account, peer)
	}
}

func (m *mockSessionEventHandler) Disconnected() {
	if m.disconnected != nil {
		m.disconnected()
	}
}

func (m *mockSessionEventHandler) RegisterCallback(title, instructions string, fields []interface{}) error {
	if m.registerCallback != nil {
		return m.registerCallback(title, instructions, fields)
	}
	return nil
}

type mockConnIOReaderWriter struct {
	read        []byte
	readIndex   int
	write       []byte
	errCount    int
	calledClose int
	err         error
}

func (iom *mockConnIOReaderWriter) Read(p []byte) (n int, err error) {
	if iom.readIndex >= len(iom.read) {
		return 0, io.EOF
	}
	i := copy(p, iom.read[iom.readIndex:])
	iom.readIndex += i
	var e error
	if iom.errCount == 0 {
		e = iom.err
	}
	iom.errCount--
	return i, e
}

func (iom *mockConnIOReaderWriter) Write(p []byte) (n int, err error) {
	iom.write = append(iom.write, p...)
	var e error
	if iom.errCount == 0 {
		e = iom.err
	}
	iom.errCount--
	return len(p), e
}

func (iom *mockConnIOReaderWriter) Close() error {
	iom.calledClose++
	return nil
}

func captureLogsEvents(c <-chan interface{}) (ret []events.Log) {
	for {
		select {
		case ev := <-c:
			switch t := ev.(type) {
			case events.Log:
				ret = append(ret, t)
			default:
				//ignore
			}
		case <-time.After(1 * time.Millisecond):
			return
		}
	}

	return
}

func assertLogContains(c *gocheck.C, ch <-chan interface{}, exp events.Log) {
	logs := captureLogsEvents(ch)

	for _, l := range logs {
		if reflect.DeepEqual(l, exp) {
			return
		}
	}

	c.Errorf("Could not finr %#v in %#v", exp, logs)
}
