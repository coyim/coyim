// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"errors"
	"fmt"
	"net"
	"strings"

	ourNet "github.com/twstrike/coyim/net"
	"golang.org/x/net/proxy"
)

// Resolve performs a DNS SRV lookup for the XMPP server that serves the given
// domain.
func Resolve(domain string) (hostport []string, err error) {
	return massage(net.LookupSRV("xmpp-client", "tcp", domain))
}

// ResolveProxy performs a DNS SRV lookup for the xmpp server that serves the given domain over the given proxy
func ResolveProxy(proxy proxy.Dialer, domain string) (hostport []string, err error) {
	return massage(ourNet.LookupSRV(proxy, "xmpp-client", "tcp", domain))
}

func massage(cname string, addrs []*net.SRV, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}

	// https://xmpp.org/rfcs/rfc6120.html#tcp-resolution-prefer
	if len(addrs) == 1 && addrs[0].Target == "." {
		return nil, errors.New("the service is decidedly not available at this domain")
	}

	ret := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		ret = append(ret, fmt.Sprintf("%s:%d",
			strings.TrimSuffix(addr.Target, "."),
			addr.Port,
		))
	}

	return ret, nil
}
