package otrclient

import (
	"fmt"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"
)

// ConversationBuilder takes a JID and returns a new conversation
type ConversationBuilder func(jid.Any) *otr3.Conversation

// OnEventHandlerCreation is a callback for any kind of event
type OnEventHandlerCreation func(jid.Any, *EventHandler, chan string, chan int)

// Sender represents an entity capable of sending messages to peers
type Sender interface {
	// Send sends a message to a peer
	Send(peer jid.Any, msg string, otr bool) error
}

// ConversationManager represents an entity capable of managing Conversations
type ConversationManager interface {
	// GetConversationWith returns the conversation for the given peer, and
	// whether the Conversation exists
	GetConversationWith(peer jid.Any) (Conversation, bool)

	// GetConversationWith returns the conversation for the given peer, and
	// creates the conversation if none exists. Additionally, returns whether the
	// conversation was created.
	EnsureConversationWith(peer jid.Any) (Conversation, bool)

	// TerminateAll terminates all existing conversations
	TerminateAll()
}

type conversationManager struct {
	// conversations maps from either a bare JID or a full to Conversation
	// the mapping here will be changed when a conversation is locked
	conversations map[string]*conversation
	sync.RWMutex

	builder    ConversationBuilder
	sender     Sender
	onCreateEH OnEventHandlerCreation

	account string
	log     coylog.Logger
}

// NewConversationManager returns a new ConversationManager
func NewConversationManager(builder ConversationBuilder, sender Sender, account string, onCreateEH OnEventHandlerCreation, l coylog.Logger) ConversationManager {
	return &conversationManager{
		conversations: make(map[string]*conversation),
		builder:       builder,
		sender:        sender,
		onCreateEH:    onCreateEH,
		account:       account,
		log:           l,
	}
}

func (m *conversationManager) getConversationWithUnlocked(peer jid.Any) (Conversation, bool) {
	pwr, pwor := jid.WithAndWithout(peer)
	if pwr != nil {
		c, ok := m.conversations[pwr.String()]
		if ok {
			return c, true
		}
	}

	c, ok := m.conversations[pwor.String()]
	if ok && pwr != nil {
		if c.locked {
			fmt.Printf("UNEXPECTED SITUATION OCCURRED - conversation with %s already locked to %s without saved correct - this shouldn't be possible\n", pwor, pwr)
		}
		c.locked = true
		c.peer = pwr
		c.eh.peer = pwr
		delete(m.conversations, pwor.String())
		m.conversations[pwr.String()] = c
	}

	return c, ok
}

func (m *conversationManager) GetConversationWith(peer jid.Any) (Conversation, bool) {
	m.Lock()
	defer m.Unlock()
	return m.getConversationWithUnlocked(peer)
}

func (m *conversationManager) EnsureConversationWith(peer jid.Any) (Conversation, bool) {
	m.Lock()
	defer m.Unlock()

	c1, ok := m.getConversationWithUnlocked(peer)
	if ok {
		return c1, false
	}

	_, locked := peer.(jid.WithResource)

	c := &conversation{
		peer:         peer,
		locked:       locked,
		s:            m.sender,
		Conversation: m.builder(peer),
	}
	c.eh = m.createEventHandler(peer, c.Conversation)

	m.conversations[peer.String()] = c

	return c, true
}

func (m *conversationManager) TerminateAll() {
	m.RLock()
	defer m.RUnlock()

	for _, c := range m.conversations {
		_ = c.EndEncryptedChat()
	}
}

func (m *conversationManager) createEventHandler(peer jid.Any, conversation *otr3.Conversation) *EventHandler {
	notificationsChan := make(chan string)
	delayedChan := make(chan int)
	eh := &EventHandler{
		delays:             make(map[int]bool),
		peer:               peer,
		account:            m.account,
		notifications:      notificationsChan,
		delayedMessageSent: delayedChan,
		log:                m.log,
	}
	m.onCreateEH(peer, eh, notificationsChan, delayedChan)
	conversation.SetSMPEventHandler(eh)
	conversation.SetErrorMessageHandler(eh)
	conversation.SetMessageEventHandler(eh)
	conversation.SetSecurityEventHandler(eh)
	return eh
}
