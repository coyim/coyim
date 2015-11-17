// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"crypto/tls"
	"io"
	"io/ioutil"
)

// Config contains options for an XMPP connection.
type Config struct {
	// InLog is an optional Writer which receives the raw contents of the
	// XML from the server.
	InLog io.Writer
	// OutLog is an optional Writer which receives the raw XML sent to the
	// server.
	OutLog io.Writer
	// Log is an optional Writer which receives human readable log messages
	// during the connection.
	Log io.Writer
	// CreateCallback, if not nil, causes a new account to be created on
	// the server. The callback is needed in order to be able to handle
	// XMPP forms.
	CreateCallback FormCallback
	// TrustedAddress, if true, means that the address passed to Dial is
	// trusted and that certificates for that name should be accepted.
	TrustedAddress bool
	// Archive determines whether we disable archiving for messages. If
	// false, XML is sent with each message to disable recording on the
	// server.
	Archive bool
	// ServerCertificateSHA256 contains the SHA-256 hash of the server's
	// leaf certificate, or may be empty to use normal X.509 verification.
	// If this is specified then normal X.509 verification is disabled.
	ServerCertificateSHA256 []byte
	// SkipTLS, if true, causes the TLS handshake to be skipped.
	// WARNING: this should only be used if Conn is already secure.
	SkipTLS bool
	// TLSConfig contains the configuration to be used by the TLS
	// handshake. If nil, sensible defaults will be used.
	TLSConfig *tls.Config
	// SkipSRVLookup skips SRV lookup during resolution of fully qualified domain
	// names. RFC 6210 section 3.2.3 recomends to skip the SRV lookup when the
	// initiating entity has a hardcoded FQDN associated with the origin domain.
	SkipSRVLookup bool
}

func (c *Config) getLog() io.Writer {
	if c.Log == nil {
		return ioutil.Discard
	}

	return c.Log
}
