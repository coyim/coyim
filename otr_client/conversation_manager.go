package otr_client

import (
	"sync"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"
)

type ConversationBuilder func(jid.Any) *otr3.Conversation
type OnEventHandlerCreation func(jid.Any, *EventHandler, chan string, chan int)

// Sender represents an entity capable of sending messages to peers
type Sender interface {
	// Send sends a message to a peer
	Send(peer jid.WithoutResource, resource jid.Resource, msg string) error
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

	// LockJID Not sure about this one - we'll have to try it
	//	LockJID(peer jid.WithResource)
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
}

// NewConversationManager returns a new ConversationManager
func NewConversationManager(builder ConversationBuilder, sender Sender, account string, onCreateEH OnEventHandlerCreation) ConversationManager {
	return &conversationManager{
		conversations: make(map[string]*conversation),
		builder:       builder,
		sender:        sender,
		onCreateEH:    onCreateEH,
		account:       account,
	}
}

func peerWithAndWithout(peer jid.Any) (jid.WithResource, jid.WithoutResource) {
	if pwr, ok := peer.(jid.WithResource); ok {
		return pwr, peer.NoResource()
	}
	return nil, peer.NoResource()
}

func (m *conversationManager) getConversationWithUnlocked(peer jid.Any) (Conversation, bool) {
	pwr, pwor := peerWithAndWithout(peer)
	if pwr != nil {
		c, ok := m.conversations[pwr.String()]
		if ok {
			return c, true
		}
	}

	c, ok := m.conversations[pwor.String()]
	return c, ok
}

// Should this lock the jid automatically? Make it manual for now...

func (m *conversationManager) GetConversationWith(peer jid.Any) (Conversation, bool) {
	m.RLock()
	defer m.RUnlock()
	return m.getConversationWithUnlocked(peer)
}

// Should this lock the jid automatically? Make it manual for now...

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
		c.EndEncryptedChat()
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
	}
	m.onCreateEH(peer, eh, notificationsChan, delayedChan)
	conversation.SetSMPEventHandler(eh)
	conversation.SetErrorMessageHandler(eh)
	conversation.SetMessageEventHandler(eh)
	conversation.SetSecurityEventHandler(eh)
	return eh
}
