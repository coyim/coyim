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
	Host   *string
	Port   *string
	Path   *string
}

var proxyTypes = [][]string{
	[]string{"tor-auto", "Automatic Tor"},
	[]string{"socks5", "SOCKS5"},
	[]string{"socks5+unix", "SOCKS5 over Unix Domain Socket"},
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
	if p.Host != "" {
		var potHost, potPort string
		potHost, potPort, err = net.SplitHostPort(p.Host)
		if err != nil && err.(*net.AddrError).Err == "missing port in address" {
			prox.Host = &p.Host
		} else {
			prox.Host = &potHost
			prox.Port = &potPort
		}
	}

	if p.Path != "" {
		prox.Path = &p.Path
	}

	return prox
}

func orEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ForProcessing represents a string valid for computer processing
func (p Proxy) ForProcessing() string {
	us := ""
	ps := ""
	compose := ""
	if p.User != nil {
		us = *p.User
		compose = "@"
		if p.Pass != nil {
			ps = ":" + *p.Pass
		}
	}

	pr := ""
	if p.Port != nil {
		pr = ":" + *p.Port
	}

	return fmt.Sprintf("%s://%s%s%s%s%s%s", p.Scheme, us, ps, compose, orEmpty(p.Host), pr, orEmpty(p.Path))
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

	return fmt.Sprintf("%s://%s%s%s%s%s%s", p.Scheme, us, ps, compose, orEmpty(p.Host), pr, orEmpty(p.Path))
}
