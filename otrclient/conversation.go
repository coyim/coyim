package otrclient

import (
	"bytes"
	"crypto/rand"
	"math"
	"math/big"

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
	GetAndWipeLastExtraKey() (usage uint32, usageData []byte, symkey []byte)

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

	lastExtraKeyUsage     uint32
	lastExtraKeyUsageData []byte
	lastExtraKeySymkey    []byte
}

func (c *conversation) StartEncryptedChat() error {
	return c.s.Send(c.peer, string(c.QueryMessage()), true)
}

func (c *conversation) EventHandler() *EventHandler {
	return c.eh
}

func isEncrypted(m otr3.ValidMessage) bool {
	return bytes.HasPrefix(m, []byte("?OTR"))
}

func (c *conversation) sendAll(toSend []otr3.ValidMessage) error {
	for _, msg := range toSend {
		isEnc := isEncrypted(msg)
		err := c.s.Send(c.peer, string(msg), isEnc)
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
	nBig, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	trace = int(nBig.Int64())

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

func (c *conversation) ReceivedSymmetricKey(usage uint32, usageData []byte, symkey []byte) {
	c.lastExtraKeyUsage = usage
	c.lastExtraKeyUsageData = usageData
	c.lastExtraKeySymkey = symkey
}

func (c *conversation) GetAndWipeLastExtraKey() (usage uint32, usageData []byte, symkey []byte) {
	usage = c.lastExtraKeyUsage
	usageData = c.lastExtraKeyUsageData
	symkey = c.lastExtraKeySymkey

	c.lastExtraKeyUsage = 0
	c.lastExtraKeyUsageData = nil
	c.lastExtraKeySymkey = nil

	return
}
