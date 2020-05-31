// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"errors"
	"net"
	"strconv"
	"strings"

	ourNet "github.com/coyim/coyim/net"
	"golang.org/x/net/proxy"
)

var (
	// ErrServiceNotAvailable means that the service is decidedly not available at this domain
	ErrServiceNotAvailable = errors.New("service not available")
)

// Resolve performs a DNS SRV lookup for the XMPP server that serves the given
// domain.
func Resolve(domain string) (hostporttls []string, hostport []string, err error) {
	_, addrs, err := net.LookupSRV("xmpps-client", "tcp", domain)
	if err != nil {
		return nil, nil, err
	}
	hostporttls, err = massage(addrs)
	if err != nil {
		return nil, nil, err
	}

	_, addrs, err = net.LookupSRV("xmpp-client", "tcp", domain)
	if err != nil {
		return nil, nil, err
	}
	hostport, err = massage(addrs)
	if err != nil {
		return nil, nil, err
	}

	return hostporttls, hostport, nil
}

// ResolveSRVWithProxy performs a DNS SRV lookup for the xmpp server that serves the given domain over the given proxy
func ResolveSRVWithProxy(proxy proxy.Dialer, domain string) (hostporttls []string, hostport []string, err error) {
	_, addrs, err := ourNet.LookupSRV(proxy, "xmpps-client", "tcp", domain)
	if err != nil {
		return nil, nil, err
	}
	hostporttls, err = massage(addrs)
	if err != nil {
		return nil, nil, err
	}

	_, addrs, err = ourNet.LookupSRV(proxy, "xmpp-client", "tcp", domain)
	if err != nil {
		return nil, nil, err
	}
	hostport, err = massage(addrs)
	if err != nil {
		return nil, nil, err
	}

	return hostporttls, hostport, nil
}

func massage(addrs []*net.SRV) ([]string, error) {
	// https://xmpp.org/rfcs/rfc6120.html#tcp-resolution-prefer
	if len(addrs) == 1 && addrs[0].Target == "." {
		return nil, ErrServiceNotAvailable
	}

	ret := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		hostport := net.JoinHostPort(
			strings.TrimSuffix(addr.Target, "."),
			strconv.Itoa(int(addr.Port)),
		)

		ret = append(ret, hostport)
	}

	return ret, nil
}
