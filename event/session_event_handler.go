package event

import (
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
)

type SessionEventHandler interface {
	Info(string)
	Warn(string)
	Alert(string)
	RosterReceived([]xmpp.RosterEntry)
	IQReceived(uid string)
	NewOTRKeys(from string, conversation *otr3.Conversation)
	OTREnded(uid string)
	MessageReceived(from, timestamp string, encrypted bool, message []byte)
	ProcessPresence(stanza *xmpp.ClientPresence, ignore, gone bool)
	Disconnected()
	RegisterCallback() xmpp.FormCallback
}
