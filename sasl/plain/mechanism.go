// Package plain implements the Plain Simple Authentication Mechanism
// according to RFC 4616
package plain

import "github.com/coyim/coyim/sasl"

var (
	// Mechanism is the PLAIN SASL mechanism
	Mechanism sasl.Mechanism = &plainMechanism{}
)

const (
	// Name is the authentication type associated with the SASL mechanism
	Name = "PLAIN"
)

// Register the SASL mechanism
func Register() {
	sasl.RegisterMechanism(Name, Mechanism)
}

type plainMechanism struct{}

func (m *plainMechanism) NewClient() sasl.Session {
	return &plain{
		state: replyChallenge{},
	}
}

type plain struct {
	user     string
	password string
	state
}

//For compatibility with GSASL, this should be called by a configurable callback
func (p *plain) SetProperty(prop sasl.Property, v string) error {
	switch prop {
	case sasl.AuthID:
		p.user = v
	case sasl.Password:
		p.password = v
	default:
		return sasl.ErrUnsupportedProperty
	}

	return nil
}

func (p *plain) Step(sasl.Token) (t sasl.Token, err error) {
	p.state, t, err = p.state.challenge(p.user, p.password)
	return
}

func (p *plain) NeedsMore() bool {
	return p.state != finished{}
}

func (p *plain) SetChannelBinding(v []byte) {
}
