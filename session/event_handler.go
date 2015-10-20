package session

import "github.com/twstrike/otr3"

type SessionEventHandler interface {
	Debug(string)
	Info(string)
	Warn(string)
	Alert(string)
	RosterReceived(*Session)
	IQReceived(uid string)
	NewOTRKeys(from string, conversation *otr3.Conversation)
	OTREnded(uid string)
	MessageReceived(s *Session, from, timestamp string, encrypted bool, message []byte)
	ProcessPresence(from, to, show, status string, gone bool)
	SubscriptionRequest(s *Session, uid string)
	Subscribed(account, peer string)
	Unsubscribe(account, peer string)
	Disconnected()
	RegisterCallback(title, instructions string, fields []interface{}) error
}
