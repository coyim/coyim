// Package scram implements the Salted Challenge Response Authentication Mechanism
// according to RFC 5802.
package scram

import "github.com/coyim/coyim/sasl"

// Register the SASL mechanisms
func Register() {
	sasl.RegisterMechanism(sha1Name, sha1Mechanism)
}
