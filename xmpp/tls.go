// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	gotls "crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp/interfaces"
)

func certName(cert *x509.Certificate) string {
    name := cert.Subject
    var b strings.Builder
    for _, org := range name.Organization {
        b.WriteString("O=" + org + "/")
    }
    for _, ou := range name.OrganizationalUnit {
        b.WriteString("OU=" + ou + "/")
    }
    if len(name.CommonName) > 0 {
        b.WriteString("CN=" + name.CommonName + "/")
    }
    return b.String()
}

// GetCipherSuiteName returns a human readable string of the cipher suite used in the state
func GetCipherSuiteName(tlsState gotls.ConnectionState) string {
	cipherSuite, ok := tlsCipherSuiteNames[tlsState.CipherSuite]
	if !ok {
		return "unknown"
	}
	return cipherSuite
}

// GetTLSVersion returns a human readable string of the TLS version used in the state
func GetTLSVersion(tlsState gotls.ConnectionState) string {
	version, ok := tlsVersionStrings[tlsState.Version]
	if !ok {
		return "unknown"
	}

	return version
}

func printTLSDetails(ll coylog.Logger, tlsState gotls.ConnectionState) {
	ll.WithField("version", GetTLSVersion(tlsState)).Info("  SSL/TLS/version")
	ll.WithField("cipherSuite", GetCipherSuiteName(tlsState)).Info("  Cipher suite")
}

// setChannelBindingData sets channel binding data on the connection based on TLS version
func setChannelBindingData(c interfaces.Conn, tlsConn tls.Conn, tlsState gotls.ConnectionState, log coylog.Logger) {
	var channelBinding []byte
	var channelBindingType string

	if tlsState.Version == gotls.VersionTLS13 {
		// For TLS 1.3, use the exporter mechanism (RFC 9266)
		// ExportKeyingMaterial is on ConnectionState, not Conn
		var err error
		channelBinding, err = tlsState.ExportKeyingMaterial("EXPORTER-Channel-Binding", nil, 32)
		if err != nil {
			log.WithError(err).Warn("Failed to export keying material for TLS 1.3 channel binding")
			channelBinding = nil
			channelBindingType = ""
		} else {
			channelBindingType = "tls-exporter"
			log.Info("Using TLS 1.3 channel binding (tls-exporter)")
		}
	} else {
		// For TLS 1.2 and earlier, use tls-unique
		channelBinding = tlsState.TLSUnique
		if channelBinding != nil {
			channelBindingType = "tls-unique"
			log.Info("Using TLS 1.2 channel binding (tls-unique)")
		}
	}

	c.SetChannelBinding(channelBinding)
	// Set channel binding type if we can
	if cc, ok := c.(*conn); ok {
		cc.SetChannelBindingType(channelBindingType)
	}
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
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	return c.SendInitialStreamHeader()
}

func (d *dialer) startTLS(c interfaces.Conn, conn net.Conn) error {
	fmt.Fprintf(c.Out(), "<starttls xmlns='%s'/>", NsTLS)

	proceed, err := nextStart(c.In(), d.log)
	if err != nil {
		return fmt.Errorf("failed to get start stanza: %w", err)
	}

	if proceed.Name.Space != NsTLS || proceed.Name.Local != "proceed" {
		return errors.New("xmpp: expected <proceed> after <starttls> but got <" + proceed.Name.Local + "> in " + proceed.Name.Space)
	}

	return d.startRawTLS(c, conn)
}

func (d *dialer) startRawTLS(c interfaces.Conn, conn net.Conn) error {
	d.log.Info("Starting TLS handshake")

	tlsConfig := c.Config().TLSConfig
	if tlsConfig == nil {
		tlsConfig = &gotls.Config{}
	}

	tlsConfig.ServerName = c.OriginDomain()
	if d.sendALPN {
		tlsConfig.NextProtos = []string{"xmpp-client"}
	}
	tlsConfig.InsecureSkipVerify = true

	tlsConn := d.tlsConnFactory(conn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		return err
	}

	tlsState := tlsConn.ConnectionState()
	printTLSDetails(d.log, tlsState)

	if err := d.verifier.Verify(tlsState, tlsConfig, c.OriginDomain()); err != nil {
		return err
	}

	// Set channel binding data based on TLS version
	// TLS 1.3 uses tls-exporter (RFC 9266), earlier versions use tls-unique
	setChannelBindingData(c, tlsConn, tlsState, d.log)
	d.bindTransport(c, tlsConn)

	return nil
}
