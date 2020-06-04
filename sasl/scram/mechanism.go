// Package scram implements the Salted Challenge Response Authentication Mechanism
// according to RFC 5802 and RFC 7677.
package scram

import (
	/* #nosec G505 */
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"github.com/coyim/coyim/sasl"
)

var (
	/* #nosec G401 */
	sha1Mechanism sasl.Mechanism = &scramMechanism{sha1.New, sha1.Size, false, true}
	/* #nosec G401 */
	sha1PlusMechanism   sasl.Mechanism = &scramMechanism{sha1.New, sha1.Size, true, true}
	sha256Mechanism     sasl.Mechanism = &scramMechanism{sha256.New, sha256.Size, false, true}
	sha256PlusMechanism sasl.Mechanism = &scramMechanism{sha256.New, sha256.Size, true, true}
	sha512Mechanism     sasl.Mechanism = &scramMechanism{sha512.New, sha512.Size, false, true}
	sha512PlusMechanism sasl.Mechanism = &scramMechanism{sha512.New, sha512.Size, true, true}
)

const (
	// Name is the authentication type associated with the SASL mechanism
	sha1Name       = "SCRAM-SHA-1"
	sha1PlusName   = "SCRAM-SHA-1-PLUS"
	sha256Name     = "SCRAM-SHA-256"
	sha256PlusName = "SCRAM-SHA-256-PLUS"
	sha512Name     = "SCRAM-SHA-512"
	sha512PlusName = "SCRAM-SHA-512-PLUS"
)

type scramMechanism struct {
	hash                  func() hash.Hash
	hashSize              int
	plus                  bool
	supportChannelBinding bool
}

func (m *scramMechanism) NewClient() sasl.Session {
	return &scram{
		state: start{m.hash, m.hashSize, m.plus, m.supportChannelBinding},
		props: make(sasl.Properties),
	}
}

type scram struct {
	state
	channelBinding []byte
	props          sasl.Properties
}

func (p *scram) SetProperty(prop sasl.Property, v string) error {
	p.props[prop] = v
	return nil
}

func (p *scram) Step(t sasl.Token) (ret sasl.Token, err error) {
	pairs := sasl.ParseAttributeValuePairs(t)
	p.state, ret, err = p.state.next(t, p.props, pairs, p.channelBinding)
	return
}

func (p *scram) NeedsMore() bool {
	return !p.state.finished()
}

func (p *scram) SetChannelBinding(v []byte) {
	p.channelBinding = v
}
