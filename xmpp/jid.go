package xmpp

import (
	"strings"
)

//JID represents a Jabber Identifier as specified in RFC 7622
type JID struct {
	LocalPart, DomainPart, ResourcePart string
}

//ParseJID parses a JID according to RFC 7622.
func ParseJID(jid string) *JID {
	domainPartBegin := strings.IndexRune(jid, '@')
	domainPartEnd := strings.IndexRune(jid, '/')

	localPart := ""
	if domainPartBegin != -1 {
		localPart = jid[0:domainPartBegin]
	}

	resourcePart := ""
	if domainPartEnd != -1 {
		resourcePart = jid[domainPartEnd+1:]
	} else {
		domainPartEnd = len(jid)
	}

	return &JID{
		LocalPart:    localPart,
		DomainPart:   jid[domainPartBegin+1 : domainPartEnd],
		ResourcePart: resourcePart,
	}
}
