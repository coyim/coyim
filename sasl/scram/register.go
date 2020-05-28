// Package scram implements the Salted Challenge Response Authentication Mechanism
// according to RFC 5802.
package scram

import "github.com/coyim/coyim/sasl"

// Register the SASL mechanisms
func Register() {
	sasl.RegisterMechanism(sha1Name, sha1Mechanism)
	sasl.RegisterMechanism(sha256Name, sha256Mechanism)
	sasl.RegisterMechanism(sha1PlusName, sha1PlusMechanism)
	sasl.RegisterMechanism(sha256PlusName, sha256PlusMechanism)
}
