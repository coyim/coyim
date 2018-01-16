package otr_client

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"
)

// This file contains primarily the type EventHandlers
// This type is used to map all active OTR event handlers
// It's done this way to simplify the logic around how Peers and JIDs are mapped to event handlers
// Initially, this will be a pure refactoring of what's in Session

type EventHandlers struct {
	handlers map[string]*EventHandler
	account  string
	onCreate func(jid.Any, *EventHandler, chan string, chan int)
}

func NewEventHandlers(account string, onCreate func(jid.Any, *EventHandler, chan string, chan int)) *EventHandlers {
	return &EventHandlers{
		handlers: make(map[string]*EventHandler),
		account:  account,
		onCreate: onCreate,
	}
}

func (ehs *EventHandlers) create(peer jid.Any, conversation *otr3.Conversation) {
	notificationsChan := make(chan string)
	delayedChan := make(chan int)
	eh := &EventHandler{
		Delays:             make(map[int]bool),
		Peer:               peer,
		Account:            ehs.account,
		Notifications:      notificationsChan,
		DelayedMessageSent: delayedChan,
	}
	ehs.onCreate(peer, eh, notificationsChan, delayedChan)
	conversation.SetSMPEventHandler(eh)
	conversation.SetErrorMessageHandler(eh)
	conversation.SetMessageEventHandler(eh)
	conversation.SetSecurityEventHandler(eh)
	ehs.Add(peer, eh)
}

func (ehs *EventHandlers) EnsureExists(peer jid.Any, conversation *otr3.Conversation) {
	_, ok := ehs.handlers[peer.String()]
	if !ok {
		ehs.create(peer, conversation)
	}
}

func (ehs *EventHandlers) Get(peer jid.Any) *EventHandler {
	return ehs.handlers[peer.String()]
}

func (ehs *EventHandlers) Add(peer jid.Any, eh *EventHandler) {
	ehs.handlers[peer.String()] = eh
}
