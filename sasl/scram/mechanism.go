// Package scram implements the Salted Challenge Response Authentication Mechanism
// according to RFC 5802.
package scram

import "../../sasl"

var (
	// Mechanism is the SCRAM-SHA1 SASL mechanism
	Mechanism sasl.Mechanism = &scramSHA1Mechanism{}
)

const (
	// Name is the authentication type associated with the SASL mechanism
	Name = "SCRAM-SHA-1"
)

// Register the SASL mechanism
func Register() {
	sasl.RegisterMechanism(Name, Mechanism)
}

type scramSHA1Mechanism struct{}

func (m *scramSHA1Mechanism) NewClient() sasl.Session {
	return &scramSHA1{
		state: firstMessage{},
		props: make(sasl.Properties),
	}
}

type scramSHA1 struct {
	state
	props sasl.Properties
}

func (p *scramSHA1) SetProperty(prop sasl.Property, v string) error {
	p.props[prop] = v
	return nil
}

func (p *scramSHA1) Step(t sasl.Token) (ret sasl.Token, err error) {
	pairs := sasl.ParseAttributeValuePairs(t)
	p.state, ret, err = p.state.challenge(t, p.props, pairs)
	return
}

func (p *scramSHA1) NeedsMore() bool {
	return p.state != finished{}
}
