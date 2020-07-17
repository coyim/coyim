// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"errors"
	"net"
	"sort"
	"strconv"
	"strings"

	ourNet "github.com/coyim/coyim/net"
	"golang.org/x/net/proxy"
)

var (
	// ErrServiceNotAvailable means that the service is decidedly not available at this domain
	ErrServiceNotAvailable = errors.New("service not available")
)

type connectEntry struct {
	host     string
	port     int
	priority int
	weight   int
	tls      bool
}

func intoConnectEntry(s string) *connectEntry {
	host, port, e := net.SplitHostPort(s)
	if e != nil {
		return nil
	}
	pp, _ := strconv.Atoi(port)
	return &connectEntry{
		host: host,
		port: pp,
		tls:  false,
	}
}

type byPriorityWeight []*connectEntry

func (s byPriorityWeight) Len() int { return len(s) }
func (s byPriorityWeight) Less(i, j int) bool {
	if s[i].priority == s[j].priority {
		return s[i].weight > s[j].weight
	}

	return s[i].priority < s[j].priority
}
func (s byPriorityWeight) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// resolve performs a DNS SRV lookup for the XMPP server that serves the given
// domain.
func resolve(domain string) (hosts []*connectEntry, err error) {
	return resolveWithCustom(domain, net.LookupSRV)
}

// resolveSRVWithProxy performs a DNS SRV lookup for the xmpp server that serves the given domain over the given proxy
func resolveSRVWithProxy(proxy proxy.Dialer, domain string) (hosts []*connectEntry, err error) {
	return resolveWithCustom(domain, resolverWithProxy(proxy))
}

type resolver func(string, string, string) (string, []*net.SRV, error)

func resolverWithProxy(p proxy.Dialer) resolver {
	return func(serv string, proto string, domain string) (string, []*net.SRV, error) {
		return ourNet.LookupSRV(p, serv, proto, domain)
	}
}

func resolveOne(part, domain string, fn resolver, tls bool) (hosts []*connectEntry, err error) {
	_, addrs, err := fn(part, "tcp", domain)
	if err != nil {
		return nil, err
	}
	hostport, err := massage(addrs, tls)
	if err != nil {
		return nil, err
	}
	return hostport, nil
}

func resolveWithCustom(domain string, fn resolver) (hosts []*connectEntry, err error) {
	hostporttls, err1 := resolveOne("xmpps-client", domain, fn, true)
	hostport, err2 := resolveOne("xmpp-client", domain, fn, false)

	result := append(hostporttls, hostport...)
	if len(result) == 0 {
		if err1 != nil {
			return nil, err1
		}
		if err2 != nil {
			return nil, err2
		}
	}

	sort.Sort(byPriorityWeight(result))

	return result, nil
}

func (c *connectEntry) String() string {
	return net.JoinHostPort(c.host, strconv.Itoa(c.port))
}

func massage(addrs []*net.SRV, tls bool) ([]*connectEntry, error) {
	// https://xmpp.org/rfcs/rfc6120.html#tcp-resolution-prefer
	if len(addrs) == 1 && addrs[0].Target == "." {
		return nil, ErrServiceNotAvailable
	}

	ret := make([]*connectEntry, 0, len(addrs))
	for _, addr := range addrs {
		hostport := &connectEntry{
			host:     strings.TrimSuffix(addr.Target, "."),
			port:     int(addr.Port),
			priority: int(addr.Priority),
			weight:   int(addr.Weight),
			tls:      tls,
		}

		ret = append(ret, hostport)
	}

	return ret, nil
}
