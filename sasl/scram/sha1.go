// Package scram implements the Salted Challenge Response Authentication Mechanism
// according to RFC 5802.
package scram

import (
	"crypto/sha1"
	"crypto/sha256"
	"hash"

	"github.com/coyim/coyim/sasl"
)

var (
	sha1Mechanism   sasl.Mechanism = &scramMechanism{sha1.New, sha1.Size, false}
	sha256Mechanism sasl.Mechanism = &scramMechanism{sha256.New, sha256.Size, false}
)

const (
	// Name is the authentication type associated with the SASL mechanism
	sha1Name   = "SCRAM-SHA-1"
	sha256Name = "SCRAM-SHA-256"
)

type scramMechanism struct {
	hash     func() hash.Hash
	hashSize int
	plus     bool
}

func (m *scramMechanism) NewClient() sasl.Session {
	return &scram{
		state: start{m.hash, m.hashSize, m.plus},
		props: make(sasl.Properties),
	}
}

type scram struct {
	state
	props sasl.Properties
}

func (p *scram) SetProperty(prop sasl.Property, v string) error {
	p.props[prop] = v
	return nil
}

func (p *scram) Step(t sasl.Token) (ret sasl.Token, err error) {
	pairs := sasl.ParseAttributeValuePairs(t)
	p.state, ret, err = p.state.next(t, p.props, pairs)
	return
}

func (p *scram) NeedsMore() bool {
	return !p.state.finished()
}
