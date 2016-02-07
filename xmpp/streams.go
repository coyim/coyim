// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"errors"
	"fmt"
)

//Send an initial stream header and receive the features required for
//continuation of the stream negotiation process.
//RFC 6120 section 4.3
func (c *Conn) sendInitialStreamHeader() error {
	if _, err := fmt.Fprintf(c.out, "<?xml version='1.0'?><stream:stream to='%s' xmlns='%s' xmlns:stream='%s' version='1.0'>\n", xmlEscape(c.originDomain), NsClient, NsStream); err != nil {
		return err
	}

	se, err := nextStart(c.in)
	if err != nil {
		return err
	}

	if se.Name.Space != NsStream || se.Name.Local != "stream" {
		//TODO: should send bad-namespace-prefix error?
		return errors.New("xmpp: expected <stream> but got <" + se.Name.Local + "> in " + se.Name.Space)
	}

	//TODO: there must be an ID in the response stream header
	//TODO: there must be an xml:lang in the response stream header
	//RFC 6120, Section 4.7.3

	// Stream features MUST follow the response stream header
	// RFC 6120, section 4.3.2
	if err := c.in.DecodeElement(&c.features, nil); err != nil {
		//TODO: should send bad-format error?
		return errors.New("xmpp: error to unmarshal <features>: " + err.Error())
	}

	return nil
}
