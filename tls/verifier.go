package tls

import "crypto/tls"

// Verifier represents something that can verify a TLS connection
type Verifier interface {
	Verify(tls.ConnectionState, tls.Config, string) error
}
