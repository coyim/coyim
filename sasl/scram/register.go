// Package scram implements the Salted Challenge Response Authentication Mechanism
// according to RFC 5802.
package scram

import "github.com/coyim/coyim/sasl"

// Register the SASL mechanisms
func Register() {
	_ = sasl.RegisterMechanism(sha1Name, sha1Mechanism)
	_ = sasl.RegisterMechanism(sha1PlusName, sha1PlusMechanism)
	_ = sasl.RegisterMechanism(sha256Name, sha256Mechanism)
	_ = sasl.RegisterMechanism(sha256PlusName, sha256PlusMechanism)
	_ = sasl.RegisterMechanism(sha512Name, sha512Mechanism)
	_ = sasl.RegisterMechanism(sha512PlusName, sha512PlusMechanism)
}
