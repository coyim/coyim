package net

import (
	"fmt"
	"net"
	"net/url"
)

// Proxy contains information about a proxy specification
type Proxy struct {
	Scheme string
	User   *string
	Pass   *string
	Host   string
	Port   *string
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

	return fmt.Sprintf("%s://%s%s%s%s%s", p.Scheme, us, ps, compose, p.Host, pr)
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
