package client

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/otr3"
	"github.com/twstrike/coyim/i18n"
)

// Conversation represents a conversation with encryption capabilities
type Conversation interface {
	Send(Sender, string, []byte) error
	Receive(Sender, string, []byte) ([]byte, error)

	StartEncryptedChat(Sender, string) error
	EndEncryptedChat(Sender, string) error

	ProvideAuthenticationSecret(Sender, string, []byte) error
	StartAuthenticate(Sender, string, string, []byte) error

	GetSSID() [8]byte
	IsEncrypted() bool
	OurFingerprint() []byte
	TheirFingerprint() []byte

	//TODO: should we expose TO and remove this from gui.ConversationWindow?
}

type conversation struct {
	to string
	*otr3.Conversation
}

func (c *conversation) StartEncryptedChat(s Sender, resource string) error {
	//TODO: review whether it should create a conversation
	//conversation, _ := m.EnsureConversationWith(peer)
	return s.Send(c.to, resource, string(c.QueryMessage())+" "+i18n.Local("Your peer has requested a private conversation with you, but your client doesn't seem to support the OTR protocol."))
}

func (c *conversation) sendAll(s Sender, resource string, toSend []otr3.ValidMessage) error {
	for _, msg := range toSend {
		err := s.Send(c.to, resource, string(msg))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *conversation) EndEncryptedChat(s Sender, resource string) error {
	toSend, err := c.End()
	if err != nil {
		return err
	}

	return c.sendAll(s, resource, toSend)
}

func (c *conversation) Send(s Sender, resource string, m []byte) error {
	toSend, err := c.Conversation.Send(m)
	if err != nil {
		return err
	}

	return c.sendAll(s, resource, toSend)
}

func (c *conversation) Receive(s Sender, resource string, m []byte) ([]byte, error) {
	plain, toSend, err := c.Conversation.Receive(m)
	err2 := c.sendAll(s, resource, toSend)

	if err != nil {
		return plain, err
	}

	return plain, err2
}

func (c *conversation) ProvideAuthenticationSecret(s Sender, resource string, m []byte) error {
	toSend, err := c.Conversation.ProvideAuthenticationSecret(m)
	if err != nil {
		return err
	}

	return c.sendAll(s, resource, toSend)
}

func (c *conversation) StartAuthenticate(s Sender, resource string, q string, m []byte) error {
	toSend, err := c.Conversation.StartAuthenticate(q, m)
	if err != nil {
		return err
	}

	return c.sendAll(s, resource, toSend)
}

func (c *conversation) OurFingerprint() []byte {
	sk := c.Conversation.GetOurCurrentKey()
	if sk == nil {
		return nil
	}

	pk := sk.PublicKey()
	if pk == nil {
		return nil
	}

	return pk.Fingerprint()
}

func (c *conversation) TheirFingerprint() []byte {
	pk := c.GetTheirKey()
	if pk == nil {
		return nil
	}

	return pk.Fingerprint()
}
