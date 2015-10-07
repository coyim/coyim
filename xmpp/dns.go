// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import "net"

// Resolve performs a DNS SRV lookup for the XMPP server that serves the given
// domain.
func Resolve(domain string) (host string, port uint16, err error) {
	_, addrs, err := net.LookupSRV("xmpp-client", "tcp", domain)

	if err != nil {
		return "", 0, err
	}

	return addrs[0].Target, addrs[0].Port, nil
}
