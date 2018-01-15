package data

import (
	"fmt"
	"strings"
)

type JIDLocal string
type JIDDomain string
type JIDResource string

type bareJID string
type fullJID string
type domainWithResource string

// JID represents any valid JID, including just a hostname, a bare jid, and a jid with a resource
type JID interface {
	Host() JIDDomain
	Representation() string
	EnsureNoResource() JIDWithoutResource
}

// JIDWithResource represents any valid JID that has a resource part
type JIDWithResource interface {
	JID
	Resource() JIDResource
	Split() (JIDWithoutResource, JIDResource)
	WithoutResource() JIDWithoutResource
}

// JIDWithoutResource represents any valid JID that does not have a resource part
type JIDWithoutResource interface {
	JID
	__forcedToNotHaveResource()
}

// BareJID represents a JID containing both a local component and a host component, but no resource component. A BareJID is a JID
type BareJID interface {
	JIDWithoutResource
	Local() JIDLocal
}

// FullJID represents a JID containing a local, host and resource component. A FullJID is a BareJID and a JID
type FullJID interface {
	JIDWithResource
	Local() JIDLocal
}

func JIDNR(s string) JIDWithoutResource {
	return ParseJID(s).EnsureNoResource()
}

func JIDR(s string) JIDWithResource {
	return ParseJID(s).(JIDWithResource)
}

func ParseJID(j string) JID {
	local := ""
	resource := ""
	left := j

	ir := strings.IndexRune(left, '/')
	if ir != -1 {
		resource = left[ir+1:]
		left = left[:ir]
	}

	ih := strings.IndexRune(left, '@')
	if ih != -1 {
		local = left[:ih]
		left = left[ih+1:]
	}

	if local != "" {
		if resource != "" {
			return fullJID(fmt.Sprintf("%s@%s/%s", local, left, resource))
		}
		return bareJID(fmt.Sprintf("%s@%s", local, left))
	}

	if resource != "" {
		return domainWithResource(fmt.Sprintf("%s/%s", left, resource))
	}

	return JIDDomain(left)
}

// Host implements JID
func (j JIDDomain) Host() JIDDomain {
	return j
}

// Representation implements JID
func (j JIDDomain) Representation() string {
	return string(j)
}

// Host implements JID
func (j bareJID) Host() JIDDomain {
	// bareJID is guaranteed to have both a local and a domain part, so that means there HAS to be an @ sign
	at := strings.IndexRune(string(j), '@')
	return JIDDomain(j[at+1:])
}

// Representation implements JID
func (j bareJID) Representation() string {
	return string(j)
}

// Local implements BareJID
func (j bareJID) Local() JIDLocal {
	// bareJID is guaranteed to have both a local and a domain part, so that means there HAS to be an @ sign
	at := strings.IndexRune(string(j), '@')
	return JIDLocal(j[:at])
}

// WithoutResource implements BareJID
func (j bareJID) WithoutResource() JIDWithoutResource {
	return j
}

// Host implements JID
func (j fullJID) Host() JIDDomain {
	return j.WithoutResource().Host()
}

// Representation implements JID
func (j fullJID) Representation() string {
	return string(j)
}

// Local implements BareJID
func (j fullJID) Local() JIDLocal {
	return j.WithoutResource().(bareJID).Local()
}

// WithoutResource implements JIDWithResource
func (j fullJID) WithoutResource() JIDWithoutResource {
	// Since a fullJID is guaranteed to contain a resource, we can assume the slash is there
	slash := strings.IndexRune(string(j), '/')
	return bareJID(j[:slash])
}

// Resource implements JIDWithResource
func (j fullJID) Resource() JIDResource {
	// Since a fullJID is guaranteed to contain a resource, we can assume the slash is there
	slash := strings.IndexRune(string(j), '/')
	return JIDResource(j[slash+1:])
}

// Split implements JIDWithResource
func (j fullJID) Split() (JIDWithoutResource, JIDResource) {
	return j.WithoutResource(), j.Resource()
}

// Host implements JID
func (j domainWithResource) Host() JIDDomain {
	return JIDDomain(j.WithoutResource().Representation())
}

// Representation implements JID
func (j domainWithResource) Representation() string {
	return string(j)
}

// Resource implements JIDWithResource
func (j domainWithResource) Resource() JIDResource {
	// Since a domainWithResource is guaranteed to contain a resource, we can assume the slash is there
	slash := strings.IndexRune(string(j), '/')
	return JIDResource(j[slash+1:])
}

// WithoutResource implements JIDWithResource
func (j domainWithResource) WithoutResource() JIDWithoutResource {
	// Since a domainWithResource is guaranteed to contain a resource, we can assume the slash is there
	slash := strings.IndexRune(string(j), '/')
	return JIDDomain(j[:slash])
}

// Split implements JIDWithResource
func (j domainWithResource) Split() (JIDWithoutResource, JIDResource) {
	return j.WithoutResource(), j.Resource()
}

func (j bareJID) __forcedToNotHaveResource()   {}
func (j JIDDomain) __forcedToNotHaveResource() {}

// EnsureNoResource implements JID
func (j JIDDomain) EnsureNoResource() JIDWithoutResource {
	return j
}

// EnsureNoResource implements JID
func (j bareJID) EnsureNoResource() JIDWithoutResource {
	return j
}

// EnsureNoResource implements JID
func (j fullJID) EnsureNoResource() JIDWithoutResource {
	return j.WithoutResource()
}

// EnsureNoResource implements JID
func (j domainWithResource) EnsureNoResource() JIDWithoutResource {
	return j.WithoutResource()
}
