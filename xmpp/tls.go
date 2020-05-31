// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/coyim/coyim/xmpp/interfaces"
)

var tlsVersionStrings = map[uint16]string{
	tls.VersionSSL30: "SSL 3.0",
	tls.VersionTLS10: "TLS 1.0",
	tls.VersionTLS11: "TLS 1.1",
	tls.VersionTLS12: "TLS 1.2",
}

var tlsCipherSuiteNames = map[uint16]string{
	tls.TLS_RSA_WITH_RC4_128_SHA:                "TLS_RSA_WITH_RC4_128_SHA",
	tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:           "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	tls.TLS_RSA_WITH_AES_128_CBC_SHA:            "TLS_RSA_WITH_AES_128_CBC_SHA",
	tls.TLS_RSA_WITH_AES_256_CBC_SHA:            "TLS_RSA_WITH_AES_256_CBC_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:        "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:          "TLS_ECDHE_RSA_WITH_RC4_128_SHA",
	tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:     "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:   "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
}

func certName(cert *x509.Certificate) string {
	name := cert.Subject
	ret := ""

	for _, org := range name.Organization {
		ret += "O=" + org + "/"
	}
	for _, ou := range name.OrganizationalUnit {
		ret += "OU=" + ou + "/"
	}
	if len(name.CommonName) > 0 {
		ret += "CN=" + name.CommonName + "/"
	}
	return ret
}

// GetCipherSuiteName returns a human readable string of the cipher suite used in the state
func GetCipherSuiteName(tlsState tls.ConnectionState) string {
	cipherSuite, ok := tlsCipherSuiteNames[tlsState.CipherSuite]
	if !ok {
		return "unknown"
	}
	return cipherSuite
}

// GetTLSVersion returns a human readable string of the TLS version used in the state
func GetTLSVersion(tlsState tls.ConnectionState) string {
	version, ok := tlsVersionStrings[tlsState.Version]
	if !ok {
		return "unknown"
	}

	return version
}

func printTLSDetails(w io.Writer, tlsState tls.ConnectionState) {
	fmt.Fprintf(w, "  SSL/TLS version: %s\n", GetTLSVersion(tlsState))
	fmt.Fprintf(w, "  Cipher suite: %s\n", GetCipherSuiteName(tlsState))
}

// RFC 6120, section 5.4
func (d *dialer) negotiateSTARTTLS(c interfaces.Conn, conn net.Conn) error {
	// RFC 6120, section 5.3
	mandatoryToNegotiate := c.Features().StartTLS.Required.Local == "required"
	if c.Config().SkipTLS && !mandatoryToNegotiate {
		return nil
	}

	// Section 5.2 states:
	// "Support for STARTTLS is REQUIRED in XMPP client and server implementations"
	if c.Features().StartTLS.XMLName.Local == "" {
		return errors.New("xmpp: server doesn't support TLS")
	}

	if err := d.startTLS(c, conn); err != nil {
		return err
	}

	return c.SendInitialStreamHeader()
}

func (d *dialer) startTLS(c interfaces.Conn, conn net.Conn) error {
	fmt.Fprintf(c.Out(), "<starttls xmlns='%s'/>", NsTLS)

	proceed, err := nextStart(c.In(), d.log)
	if err != nil {
		return err
	}

	if proceed.Name.Space != NsTLS || proceed.Name.Local != "proceed" {
		return errors.New("xmpp: expected <proceed> after <starttls> but got <" + proceed.Name.Local + "> in " + proceed.Name.Space)
	}

	l := c.Config().GetLog()
	io.WriteString(l, "Starting TLS handshake\n")

	tlsConfig := c.Config().TLSConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
	}

	tlsConfig.ServerName = c.OriginDomain()
	tlsConfig.InsecureSkipVerify = true

	tlsConn := d.tlsConnFactory(conn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		return err
	}

	tlsState := tlsConn.ConnectionState()
	printTLSDetails(l, tlsState)

	if err = d.verifier.Verify(tlsState, tlsConfig, c.OriginDomain()); err != nil {
		return err
	}

	c.SetChannelBinding(tlsState.TLSUnique)
	d.bindTransport(c, tlsConn)

	return nil
}
