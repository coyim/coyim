package otrclient

import (
	"math/rand"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"
)

// Conversation represents a conversation with encryption capabilities
type Conversation interface {
	Send([]byte) (trace int, err error)
	Receive([]byte) ([]byte, error)

	StartEncryptedChat() error
	EndEncryptedChat() error

	ProvideAuthenticationSecret([]byte) error
	StartAuthenticate(string, []byte) error
	AbortAuthentication() error

	GetSSID() [8]byte
	IsEncrypted() bool
	OurFingerprint() []byte
	TheirFingerprint() []byte

	CreateExtraSymmetricKey() ([]byte, error)

	EventHandler() *EventHandler
}

// JID considerations:
// When a conversation is first started, it can sometimes have a resource, or sometimes not
// However, from the point that we receive a message from the peer, we will ALWAYS have a resource
// At that point, the conversation will be locked in, and we will never go back to not having that resource again
// Thus, we will have a "locked" field. If it's set, we will never change the peer

type conversation struct {
	peer   jid.Any
	locked bool

	s Sender

	eh *EventHandler

	*otr3.Conversation
}

func (c *conversation) StartEncryptedChat() error {
	return c.s.Send(c.peer, string(c.QueryMessage()))
}

func (c *conversation) EventHandler() *EventHandler {
	return c.eh
}

func (c *conversation) sendAll(toSend []otr3.ValidMessage) error {
	for _, msg := range toSend {
		err := c.s.Send(c.peer, string(msg))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *conversation) EndEncryptedChat() error {
	toSend, err := c.End()
	if err != nil {
		return err
	}

	return c.sendAll(toSend)
}

func (c *conversation) Send(m []byte) (trace int, err error) {
	trace = rand.Int()
	toSend, err := c.Conversation.Send(m, trace)
	if err != nil {
		return 0, err
	}

	return trace, c.sendAll(toSend)
}

func (c *conversation) Receive(m []byte) ([]byte, error) {
	plain, toSend, err := c.Conversation.Receive(m)
	err2 := c.sendAll(toSend)

	if err != nil {
		return plain, err
	}

	return plain, err2
}

func (c *conversation) ProvideAuthenticationSecret(m []byte) error {
	toSend, err := c.Conversation.ProvideAuthenticationSecret(m)
	if err != nil {
		return err
	}

	return c.sendAll(toSend)
}

func (c *conversation) StartAuthenticate(q string, m []byte) error {
	toSend, err := c.Conversation.StartAuthenticate(q, m)
	if err != nil {
		return err
	}

	return c.sendAll(toSend)
}

func (c *conversation) AbortAuthentication() error {
	toSend, err := c.Conversation.AbortAuthentication()
	if err != nil {
		return err
	}
	return c.sendAll(toSend)
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

const usageFileTransfer = uint32(123)

func (c *conversation) CreateExtraSymmetricKey() ([]byte, error) {
	key, toSend, err := c.Conversation.UseExtraSymmetricKey(usageFileTransfer, nil)
	if err != nil {
		return nil, err
	}
	err = c.sendAll(toSend)
	return key, err
}
