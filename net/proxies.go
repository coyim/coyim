package net

import (
	"fmt"
	"net"
	"net/url"

	"github.com/twstrike/coyim/i18n"
)

// Proxy contains information about a proxy specification
type Proxy struct {
	Scheme string
	User   *string
	Pass   *string
	Host   string
	Port   *string
}

var proxyTypes = [][]string{
	[]string{"tor-auto", "Automatic Tor"},
	[]string{"socks4", "SOCKS4"},
	[]string{"socks5", "SOCKS5"},
}

// FindProxyTypeFor returns the index of the proxy type given
func FindProxyTypeFor(s string) int {
	for ix, px := range proxyTypes {
		if px[0] == s {
			return ix
		}
	}

	return -1
}

// GetProxyTypeNames will yield all i18n proxy names to the function
func GetProxyTypeNames(f func(string)) {
	for _, px := range proxyTypes {
		f(i18n.Local(px[1]))
	}
}

// GetProxyTypeFor will return the proxy type for the given i18n proxy name
func GetProxyTypeFor(act string) string {
	for _, px := range proxyTypes {
		if act == i18n.Local(px[1]) {
			return px[0]
		}
	}
	return ""
}

// ParseProxy parses the given specification and returns a Proxy object with it
func ParseProxy(px string) Proxy {
	prox := Proxy{}
	p, _ := url.Parse(px)
	prox.Scheme = p.Scheme
	if p.User != nil {
		u := p.User.Username()
		prox.User = &u
		pas, pasSet := p.User.Password()
		if pasSet {
			prox.Pass = &pas
		}
	}
	var err error
	var potPort string
	prox.Host, potPort, err = net.SplitHostPort(p.Host)
	if err != nil && err.(*net.AddrError).Err == "missing port in address" {
		prox.Host = p.Host
	} else {
		prox.Port = &potPort
	}

	return prox
}

// ForPresentation represents a string valid for user presentation - blanking out the password
func (p Proxy) ForPresentation() string {
	us := ""
	ps := ""
	compose := ""
	if p.User != nil {
		us = *p.User
		compose = "@"
		if p.Pass != nil {
			ps = ":*****"
		}
	}

	pr := ""
	if p.Port != nil {
		pr = ":" + *p.Port
	}

	return fmt.Sprintf("%s://%s%s%s%s%s", p.Scheme, us, ps, compose, p.Host, pr)
}
