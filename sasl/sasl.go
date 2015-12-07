package sasl

import (
	"errors"
	"sync"
)

//TODO: with libgsasl

//TODO: how to set Mechanism properties?
//TODO: libgsasl uses callbacks to ask for user data (client) or decide if the
//credential is authenticated (server). See: https://www.gnu.org/software/gsasl/manual/gsasl.html#Callback-Functions

// Session represents an authorization session
type Session interface {
	SetProperty(Property, string) error
	Step(Token) (Token, error)
	NeedsMore() bool
}

// Mechanism represents an SASL mechanism
type Mechanism interface {
	NewClient() Session
}

var registry = struct {
	sync.Mutex
	m map[string]Mechanism
}{
	m: make(map[string]Mechanism),
}

// RegisterMechanism registers an SASL mechanism for a name
func RegisterMechanism(name string, m Mechanism) error {
	registry.Lock()
	defer registry.Unlock()

	if _, ok := registry.m[name]; ok {
		return ErrMechanismAlreadyRegistered
	}

	registry.m[name] = m
	return nil
}

// ClientSupport returns whether there is client-side support for a specified mechanism
func ClientSupport(mechanism string) bool {
	registry.Lock()
	defer registry.Unlock()

	_, ok := registry.m[mechanism]
	return ok
}

var (
	// ErrMechanismAlreadyRegistered indicates an attempt to register a duplicate mechanism
	ErrMechanismAlreadyRegistered = errors.New("the mechanism already registered")
	// ErrUnsupportedMechanism indicates an attempt to use an unregistered mechanism
	ErrUnsupportedMechanism = errors.New("the requested mechanism is not supported")
	// ErrUnsupportedProperty indicates an attempt to set a property unsupported by a mechanism
	ErrUnsupportedProperty = errors.New("unsupported property")
)

// NewClient returns a client session for a SASL mechanism
func NewClient(mechanism string) (Session, error) {
	m, ok := registry.m[mechanism]
	if !ok {
		return nil, ErrUnsupportedMechanism
	}

	return m.NewClient(), nil
}
