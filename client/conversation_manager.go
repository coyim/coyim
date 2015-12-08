package client

import (
	"sync"

	"github.com/twstrike/otr3"
)

// Conversation represents a conversation with encryption capabilities
type Conversation interface {
	Send(otr3.ValidMessage) ([]otr3.ValidMessage, error)
	Receive(otr3.ValidMessage) (otr3.MessagePlaintext, []otr3.ValidMessage, error)
	End() ([]otr3.ValidMessage, error)

	GetSSID() [8]byte
	IsEncrypted() bool
	QueryMessage() otr3.ValidMessage

	ProvideAuthenticationSecret([]byte) ([]otr3.ValidMessage, error)
	StartAuthenticate(string, []byte) ([]otr3.ValidMessage, error)

	GetOurCurrentKey() otr3.PrivateKey
	GetTheirKey() otr3.PublicKey
}

// ConversationBuilder represents an entity capable of building Conversations
type ConversationBuilder interface {
	// NewConversation returns a new conversation to a peer
	NewConversation(peer string) Conversation
}

// Sender represents an entity capable of sending messages to peers
//TODO: this assumes there is no more than one simultaneous conversations with a given peer
type Sender interface {
	// Send sends a message to a peer
	Send(peer, msg string) error
}

// ConversationManager represents an entity capable of managing Conversations
type ConversationManager interface {
	// GetConversationWith returns the conversation for the given peer, and
	// whether the Conversation exists
	GetConversationWith(peer string) (Conversation, bool)

	// GetConversationWith returns the conversation for the given peer, and
	// creates the conversation if none exists. Additionally, returns whether the
	// conversation was created.
	EnsureConversationWith(peer string) (Conversation, bool)

	// Conversations return all conversations currently managed
	Conversations() map[string]Conversation

	// StartEncryptedChatWith starts an encrypted chat with the given peer
	StartEncryptedChatWith(peer string) error

	// TerminateConversationWith terminates a conversation with a peer
	TerminateConversationWith(peer string) error

	// TerminateAll terminates all existing conversations
	TerminateAll()
}

type conversationManager struct {
	// conversations maps from a bare JID to Conversation
	conversations map[string]Conversation
	sync.RWMutex

	builder ConversationBuilder
	sender  Sender
}

// NewConversationManager returns a new ConversationManager
func NewConversationManager(builder ConversationBuilder, sender Sender) ConversationManager {
	return &conversationManager{
		conversations: make(map[string]Conversation),
		builder:       builder,
		sender:        sender,
	}
}

func (m *conversationManager) sendMsg(peer, msg string) error {
	return m.sender.Send(peer, msg)
}

func (m *conversationManager) GetConversationWith(peer string) (Conversation, bool) {
	m.RLock()
	defer m.RUnlock()
	conversation, ok := m.conversations[peer]
	return conversation, ok
}

func (m *conversationManager) Conversations() map[string]Conversation {
	m.Lock()
	defer m.Unlock()

	ret := make(map[string]Conversation)
	for to, c := range m.conversations {
		ret[to] = c
	}

	return ret
}

func (m *conversationManager) EnsureConversationWith(peer string) (Conversation, bool) {
	if c, ok := m.GetConversationWith(peer); ok {
		return c, false
	}

	m.Lock()
	defer m.Unlock()
	conversation := m.builder.NewConversation(peer)
	m.conversations[peer] = conversation

	return conversation, true
}

func (m *conversationManager) StartEncryptedChatWith(peer string) error {
	//TODO: review whether it should create a conversation
	conversation, _ := m.EnsureConversationWith(peer)
	return m.sendMsg(peer, string(conversation.QueryMessage()))
}

func (m *conversationManager) TerminateAll() {
	m.Lock()
	defer m.Unlock()

	for peer := range m.conversations {
		m.TerminateConversationWith(peer)
	}
}

func (m *conversationManager) TerminateConversationWith(peer string) error {
	c, ok := m.GetConversationWith(peer)
	if !ok {
		return nil
	}

	msgs, err := c.End()
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		err := m.sendMsg(peer, string(msg))
		if err != nil {
			return err
		}
	}

	//TODO: add wipe for conversation
	//conversation.Wipe()

	return nil
}
