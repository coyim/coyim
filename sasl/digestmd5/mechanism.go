// Package digestmd5 implements the Digest Authentication as a SASL Mechanism
// according to RFC 2831.
package digestmd5

import "../../sasl"

var (
	// Mechanism is the DIGEST-MD5 SASL mechanism
	Mechanism sasl.Mechanism = &digestMD5Mechanism{}
)

const (
	// Name is the authentication type associated with the SASL mechanism
	Name = "DIGEST-MD5"
)

// Register the SASL mechanism
func Register() {
	sasl.RegisterMechanism(Name, Mechanism)
}

type digestMD5Mechanism struct{}

func (m *digestMD5Mechanism) NewClient() sasl.Session {
	return &digestMD5{
		digestState: clientChallenge{},
		props:       make(sasl.Properties),
	}
}

type digestMD5 struct {
	digestState
	props sasl.Properties
}

func (p *digestMD5) SetProperty(prop sasl.Property, v string) error {
	p.props[prop] = v
	return nil
}

func (p *digestMD5) Step(t sasl.Token) (ret sasl.Token, err error) {
	pairs := sasl.ParseAttributeValuePairs(t)
	p.digestState, ret, err = p.digestState.challenge(p.props, pairs)
	return
}

func (p *digestMD5) NeedsMore() bool {
	return p.digestState != finished{}
}
