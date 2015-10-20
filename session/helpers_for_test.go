package session

import (
	"io"

	"github.com/twstrike/otr3"
)

type mockSessionEventHandler struct {
	debug               func(string)
	info                func(string)
	warn                func(string)
	alert               func(string)
	rosterReceived      func(*Session)
	iqReceived          func(uid string)
	newOTRKeys          func(from string, conversation *otr3.Conversation)
	otrEnded            func(uid string)
	messageReceived     func(s *Session, from, timestamp string, encrypted bool, message []byte)
	processPresence     func(from, to, show, status string, gone bool)
	subscriptionRequest func(s *Session, uid string)
	subscribed          func(account, peer string)
	unsubscribe         func(account, peer string)
	disconnected        func()
	registerCallback    func(title, instructions string, fields []interface{}) error
}

func (m *mockSessionEventHandler) Debug(v string) {
	if m.debug != nil {
		m.debug(v)
	}
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

func (m *mockSessionEventHandler) RosterReceived(s *Session) {
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

func (m *mockSessionEventHandler) MessageReceived(s *Session, from, timestamp string, encrypted bool, message []byte) {
	if m.messageReceived != nil {
		m.messageReceived(s, from, timestamp, encrypted, message)
	}
}

func (m *mockSessionEventHandler) ProcessPresence(from, to, show, status string, gone bool) {
	if m.processPresence != nil {
		m.processPresence(from, to, show, status, gone)
	}
}

func (m *mockSessionEventHandler) SubscriptionRequest(s *Session, uid string) {
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
	read      []byte
	readIndex int
	write     []byte
	errCount  int
	err       error
}

func (in *mockConnIOReaderWriter) Read(p []byte) (n int, err error) {
	if in.readIndex >= len(in.read) {
		return 0, io.EOF
	}
	i := copy(p, in.read[in.readIndex:])
	in.readIndex += i
	var e error
	if in.errCount == 0 {
		e = in.err
	}
	in.errCount--
	return i, e
}

func (out *mockConnIOReaderWriter) Write(p []byte) (n int, err error) {
	out.write = append(out.write, p...)
	var e error
	if out.errCount == 0 {
		e = out.err
	}
	out.errCount--
	return len(p), e
}
