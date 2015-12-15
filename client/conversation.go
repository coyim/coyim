package client

import "github.com/twstrike/otr3"

// Conversation represents a conversation with encryption capabilities
type Conversation interface {
	Send(Sender, []byte) error
	Receive(Sender, []byte) ([]byte, error)

	StartEncryptedChat(Sender) error
	EndEncryptedChat(Sender) error

	ProvideAuthenticationSecret(Sender, []byte) error
	StartAuthenticate(Sender, string, []byte) error

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

func (c *conversation) StartEncryptedChat(s Sender) error {
	//TODO: review whether it should create a conversation
	//conversation, _ := m.EnsureConversationWith(peer)
	return s.Send(c.to, string(c.QueryMessage()))
}

func (c *conversation) sendAll(s Sender, toSend []otr3.ValidMessage) error {
	for _, msg := range toSend {
		err := s.Send(c.to, string(msg))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *conversation) EndEncryptedChat(s Sender) error {
	toSend, err := c.End()
	if err != nil {
		return err
	}

	return c.sendAll(s, toSend)
}

func (c *conversation) Send(s Sender, m []byte) error {
	toSend, err := c.Conversation.Send(m)
	if err != nil {
		return err
	}

	return c.sendAll(s, toSend)
}

func (c *conversation) Receive(s Sender, m []byte) ([]byte, error) {
	plain, toSend, err := c.Conversation.Receive(m)
	if err != nil {
		return nil, err
	}

	return plain, c.sendAll(s, toSend)
}

func (c *conversation) ProvideAuthenticationSecret(s Sender, m []byte) error {
	toSend, err := c.Conversation.ProvideAuthenticationSecret(m)
	if err != nil {
		return err
	}

	return c.sendAll(s, toSend)
}

func (c *conversation) StartAuthenticate(s Sender, q string, m []byte) error {
	toSend, err := c.Conversation.StartAuthenticate(q, m)
	if err != nil {
		return err
	}

	return c.sendAll(s, toSend)
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
