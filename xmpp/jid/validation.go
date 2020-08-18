package jid

import (
	"net"
	"regexp"
	"strings"
)

// ValidLocal checks whether the given string is a valid local part of a JID
func ValidLocal(s string) bool {
	l := len(s)
	if l == 0 || l > 1023 {
		return false
	}

	if strings.ContainsAny(s, "\"&'/:<>@") {
		return false
	}

	return true
}

// Patterns taken from the govalidator project
var dnsName = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
var dnsReg = regexp.MustCompile(dnsName)

// ValidDomain returns true if the given string is a valid domain part for a JID
func ValidDomain(s string) bool {
	l := len(s)
	if l == 0 || l > 1023 {
		return false
	}

	res := net.ParseIP(s)
	if res != nil {
		return true
	}

	return dnsReg.MatchString(s)
}

// ValidResource returns true if the given string is a valid resource part for a JID.
// Note that a resource part is allowed to contain / and @ characters
func ValidResource(s string) bool {
	l := len(s)
	if l == 0 || l > 1023 {
		return false
	}

	return true
}

// ValidJID returns true if the given string is any of the possible JID types
func ValidJID(s string) bool {
	return ValidFullJID(s) || ValidBareJID(s) || ValidDomain(s) || ValidDomainWithResource(s)
}

// ValidBareJID returns true if the given string is a valid bare JID. This function will true for full JIDs as well as
// bare JIDs
func ValidBareJID(s string) bool {
	// TODO
	return false
}

// ValidFullJID returns true if the given string is a valid full JID
func ValidFullJID(s string) bool {
	// TODO
	return false
}

// ValidDomainWithResource returns true if the given string a valid domain with resource part. This wil return true for
// a full JID, as well as a domain with JID
func ValidDomainWithResource(s string) bool {
	// TODO
	return false
}
