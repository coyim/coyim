// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
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

func printTLSDetails(w io.Writer, tlsState tls.ConnectionState) {
	version, ok := tlsVersionStrings[tlsState.Version]
	if !ok {
		version = "unknown"
	}

	cipherSuite, ok := tlsCipherSuiteNames[tlsState.CipherSuite]
	if !ok {
		cipherSuite = "unknown"
	}

	fmt.Fprintf(w, "  SSL/TLS version: %s\n", version)
	fmt.Fprintf(w, "  Cipher suite: %s\n", cipherSuite)
}

// RFC 6120, section 5
type tlsStartTLS struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required xml.Name `xml:"required"`
}

type tlsProceed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

type tlsFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

// RFC 6120, section 5.4
func (d *Dialer) negotiateSTARTTLS(c *Conn, conn net.Conn) error {
	// RFC 6120, section 5.3
	mandatoryToNegotiate := c.features.StartTLS.Required.Local == "required"
	if c.config.SkipTLS && !mandatoryToNegotiate {
		return nil
	}

	// Section 5.2 states:
	// "Support for STARTTLS is REQUIRED in XMPP client and server implementations"
	if c.features.StartTLS.XMLName.Local == "" {
		return errors.New("xmpp: server doesn't support TLS")
	}

	if err := d.startTLS(c, conn); err != nil {
		return err
	}

	return c.sendInitialStreamHeader()
}

func (d *Dialer) startTLS(c *Conn, conn net.Conn) error {
	address := d.GetServer()

	fmt.Fprintf(c.out, "<starttls xmlns='%s'/>", NsTLS)

	proceed, err := nextStart(c.in)
	if err != nil {
		return err
	}

	if proceed.Name.Space != NsTLS || proceed.Name.Local != "proceed" {
		return errors.New("xmpp: expected <proceed> after <starttls> but got <" + proceed.Name.Local + "> in " + proceed.Name.Space)
	}

	l := c.config.getLog()
	io.WriteString(l, "Starting TLS handshake\n")

	var tlsConfig tls.Config
	if c.config.TLSConfig != nil {
		tlsConfig = *c.config.TLSConfig
	}
	tlsConfig.ServerName = c.originDomain
	tlsConfig.InsecureSkipVerify = true

	tlsConn := tls.Client(conn, &tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		return err
	}

	tlsState := tlsConn.ConnectionState()
	printTLSDetails(l, tlsState)

	haveCertHash := len(c.config.ServerCertificateSHA256) != 0
	if haveCertHash {
		h := sha256.New()
		h.Write(tlsState.PeerCertificates[0].Raw)
		if digest := h.Sum(nil); !bytes.Equal(digest, c.config.ServerCertificateSHA256) {
			return fmt.Errorf("xmpp: server certificate does not match expected hash (got: %x, want: %x)",
				digest, c.config.ServerCertificateSHA256)
		}
	} else {
		if len(tlsState.PeerCertificates) == 0 {
			return errors.New("xmpp: server has no certificates")
		}

		opts := x509.VerifyOptions{
			Intermediates: x509.NewCertPool(),
			Roots:         tlsConfig.RootCAs,
		}
		for _, cert := range tlsState.PeerCertificates[1:] {
			opts.Intermediates.AddCert(cert)
		}
		verifiedChains, err := tlsState.PeerCertificates[0].Verify(opts)
		if err != nil {
			return errors.New("xmpp: failed to verify TLS certificate: " + err.Error())
		}

		for i, cert := range verifiedChains[0] {
			fmt.Fprintf(l, "  certificate %d: %s\n", i, certName(cert))
		}
		leafCert := verifiedChains[0][0]

		if err := leafCert.VerifyHostname(c.originDomain); err != nil {
			if c.config.TrustedAddress {
				fmt.Fprintf(l, "Certificate fails to verify against domain in username: %s\n", err)
				host, _, err := net.SplitHostPort(address)
				if err != nil {
					return errors.New("xmpp: failed to split address when checking whether TLS certificate is valid: " + err.Error())
				}

				if err = leafCert.VerifyHostname(host); err != nil {
					return errors.New("xmpp: failed to match TLS certificate to address after failing to match to username: " + err.Error())
				}
				fmt.Fprintf(l, "Certificate matches against trusted server hostname: %s\n", host)
			} else {
				return errors.New("xmpp: failed to match TLS certificate to name: " + err.Error())
			}
		}
	}

	d.bindTransport(c, tlsConn)

	return nil
}
