package sasl

import "fmt"

// Property represents a SASL mechanism property
type Property int

// All SASL mechanism properties, as defined by GSASL
const (
	AuthID Property = iota
	Password
	AuthZID
	Realm
	Service
	QOP
	ClientNonce
)

// Properties represents a map of property and value
type Properties map[Property]string

// PropertyMissingError represents an error due a missing property
type PropertyMissingError struct {
	Property
}

func (e PropertyMissingError) Error() string {
	return fmt.Sprintf("missing property %q", e.Property)
}
