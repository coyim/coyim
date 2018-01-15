package otr_client

import (
	"sync"

	"github.com/coyim/otr3"
)

// ConversationBuilder represents an entity capable of building Conversations
type ConversationBuilder interface {
	// NewConversation returns a new conversation to a peer
	NewConversation(peer string) *otr3.Conversation
}

// Sender represents an entity capable of sending messages to peers
//TODO: this assumes there is no more than one simultaneous conversations with a given peer
type Sender interface {
	// Send sends a message to a peer
	Send(peer, resource, msg string) error
}

// ConversationManager represents an entity capable of managing Conversations
type ConversationManager interface {
	// GetConversationWith returns the conversation for the given peer, and
	// whether the Conversation exists
	GetConversationWith(peer, resource string) (Conversation, bool)

	// GetConversationWith returns the conversation for the given peer, and
	// creates the conversation if none exists. Additionally, returns whether the
	// conversation was created.
	EnsureConversationWith(peer, resource string) (Conversation, bool)

	// Conversations return all conversations currently managed
	Conversations() map[string]Conversation

	// TerminateAll terminates all existing conversations
	TerminateAll()
}

type conversationManager struct {
	// conversations maps from a bare JID to Conversation
	conversations map[string]*conversation
	sync.RWMutex

	builder ConversationBuilder
	sender  Sender
}

// NewConversationManager returns a new ConversationManager
func NewConversationManager(builder ConversationBuilder, sender Sender) ConversationManager {
	return &conversationManager{
		conversations: make(map[string]*conversation),
		builder:       builder,
		sender:        sender,
	}
}

func (m *conversationManager) GetConversationWith(peer, resource string) (Conversation, bool) {
	m.RLock()
	defer m.RUnlock()
	c, ok := m.conversations[peer]
	if ok && c.resource != "" && resource != "" && c.resource != resource {
		return c, false
	}
	if ok {
		c.resource = resource
	}
	return c, ok
}

func (m *conversationManager) Conversations() map[string]Conversation {
	m.RLock()
	defer m.RUnlock()

	ret := make(map[string]Conversation)
	for to, c := range m.conversations {
		ret[to] = c
	}

	return ret
}

func (m *conversationManager) EnsureConversationWith(peer, resource string) (Conversation, bool) {
	m.Lock()
	defer m.Unlock()

	c, ok := m.conversations[peer]
	if ok && (c.resource == "" || resource == "" || c.resource == resource) {
		c.resource = resource
		return c, true
	}

	if ok {
		m.terminateConversationWith(peer, c.resource)
	}

	c = &conversation{
		to:           peer,
		resource:     resource,
		Conversation: m.builder.NewConversation(peer),
	}
	m.conversations[peer] = c

	return c, true
}

func (m *conversationManager) TerminateAll() {
	m.RLock()
	defer m.RUnlock()

	for peer := range m.conversations {
		m.terminateConversationWith(peer, "")
	}
}

func (m *conversationManager) terminateConversationWith(peer, resource string) error {
	if c, ok := m.conversations[peer]; ok {
		return c.EndEncryptedChat(m.sender, resource)
	}

	return nil
}
