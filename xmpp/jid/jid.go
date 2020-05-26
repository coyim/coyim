package jid

import (
	"fmt"
	"strings"
)

// Local represents the local part of a JID
type Local string

// Domain represents the domain part of a JID
type Domain string

// Resource represents the resource part of a JID
type Resource string

type bare string
type full string
type domainWithResource string

// Any represents any valid JID, including just a hostname, a bare jid, and a jid with a resource
type Any interface {
	// Host will always return the domain component, since all JIDs have one
	Host() Domain
	// String will return the natural string representation of the JID
	String() string
	// WithResource will return a new JID containing the resource component specified. If the JID already had a resource, it will be replaced
	WithResource(Resource) WithResource
	// MaybeWithResource will act like WithResource, if the argument is anything but a blank resource.
	// Otherwise it will return itself without a resource
	MaybeWithResource(Resource) Any
	// NoResource will ensure that the JID returned doesn't have a resource
	NoResource() WithoutResource
	// Potential resource returns the resource if one exists, or the blank resource otherwise
	PotentialResource() Resource
	// PotentialSplit will return the result of calling WithoutResource and PotentialResource
	PotentialSplit() (WithoutResource, Resource)
}

// WithResource represents any valid JID that has a resource part
type WithResource interface {
	Any
	// Resource will return the resource
	Resource() Resource
	// Split will return the JID split into the part without resource and the part with resource
	Split() (WithoutResource, Resource)
}

// WithoutResource represents any valid JID that does not have a resource part
type WithoutResource interface {
	Any
	_ForcedToNotHaveResource()
}

// Bare represents a JID containing both a local component and a host component, but no resource component. A Bare is an Any
type Bare interface {
	WithoutResource
	WithLocal
}

// Full represents a JID containing a local, host and resource component. A Full is a Bare and an Any
type Full interface {
	WithResource
	WithLocal
}

// WithLocal represents a JID that has a Local port
type WithLocal interface {
	// Local returns the local part of the JID
	Local() Local
}

// NR returns a JID without a resource
func NR(s string) WithoutResource {
	return Parse(s).NoResource()
}

// R returns a JID with resource. This method will fail if the object doesn't have a resource
func R(s string) WithResource {
	return Parse(s).(WithResource)
}

// Parse will parse the given string and return the most specific JID type that matches it
func Parse(j string) Any {
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
			return full(fmt.Sprintf("%s@%s/%s", local, left, resource))
		}
		return bare(fmt.Sprintf("%s@%s", local, left))
	}

	if resource != "" {
		return domainWithResource(fmt.Sprintf("%s/%s", left, resource))
	}

	return Domain(left)
}

// Host implements Any
func (j Domain) Host() Domain {
	return j
}

// String implements Any
func (j Domain) String() string {
	return string(j)
}

// PotentialResource implements Any
func (j Domain) PotentialResource() Resource {
	return Resource("")
}

// PotentialSplit implements Any
func (j Domain) PotentialSplit() (WithoutResource, Resource) {
	return j, j.PotentialResource()
}

// Host implements Any
func (j bare) Host() Domain {
	// bare is guaranteed to have both a local and a domain part, so that means there HAS to be an @ sign
	at := strings.IndexRune(string(j), '@')
	return Domain(j[at+1:])
}

// String implements Any
func (j bare) String() string {
	return string(j)
}

// PotentialResource implements Any
func (j bare) PotentialResource() Resource {
	return Resource("")
}

// PotentialSplit implements Any
func (j bare) PotentialSplit() (WithoutResource, Resource) {
	return j, j.PotentialResource()
}

// Local implements Bare
func (j bare) Local() Local {
	// bareJID is guaranteed to have both a local and a domain part, so that means there HAS to be an @ sign
	at := strings.IndexRune(string(j), '@')
	return Local(j[:at])
}

// NoResource implements Any
func (j bare) NoResource() WithoutResource {
	return j
}

// WithResource implements Any
func (j bare) WithResource(r Resource) WithResource {
	return R(j.String() + "/" + string(r))
}

// MaybeWithResource implements Any
func (j bare) MaybeWithResource(r Resource) Any {
	return Parse(j.String() + "/" + string(r))
}

// Host implements Any
func (j full) Host() Domain {
	return j.NoResource().Host()
}

// String implements Any
func (j full) String() string {
	return string(j)
}

// MaybeWithResource implements Any
func (j full) MaybeWithResource(r Resource) Any {
	return j.NoResource().MaybeWithResource(r)
}

// WithResource implements Any
func (j full) WithResource(r Resource) WithResource {
	return j.NoResource().WithResource(r)
}

// PotentialResource implements Any
func (j full) PotentialResource() Resource {
	return j.Resource()
}

// PotentialSplit implements Any
func (j full) PotentialSplit() (WithoutResource, Resource) {
	return j.Split()
}

// Local implements Bare
func (j full) Local() Local {
	return j.NoResource().(bare).Local()
}

// NoResource implements WithResource
func (j full) NoResource() WithoutResource {
	// Since a fullJID is guaranteed to contain a resource, we can assume the slash is there
	slash := strings.IndexRune(string(j), '/')
	return bare(j[:slash])
}

// Resource implements WithResource
func (j full) Resource() Resource {
	// Since a full is guaranteed to contain a resource, we can assume the slash is there
	slash := strings.IndexRune(string(j), '/')
	return Resource(j[slash+1:])
}

// Split implements WithResource
func (j full) Split() (WithoutResource, Resource) {
	return j.NoResource(), j.Resource()
}

// Host implements Any
func (j domainWithResource) Host() Domain {
	return Domain(j.NoResource().String())
}

// String implements Any
func (j domainWithResource) String() string {
	return string(j)
}

// MaybeWithResource implements Any
func (j domainWithResource) MaybeWithResource(r Resource) Any {
	return j.NoResource().MaybeWithResource(r)
}

// WithResource implements Any
func (j domainWithResource) WithResource(r Resource) WithResource {
	return j.NoResource().WithResource(r)
}

// PotentialResource implements Any
func (j domainWithResource) PotentialResource() Resource {
	return j.Resource()
}

// PotentialSplit implements Any
func (j domainWithResource) PotentialSplit() (WithoutResource, Resource) {
	return j.Split()
}

// Resource implements WithResource
func (j domainWithResource) Resource() Resource {
	// Since a domainWithResource is guaranteed to contain a resource, we can assume the slash is there
	slash := strings.IndexRune(string(j), '/')
	return Resource(j[slash+1:])
}

// NoResource implements WithResource
func (j domainWithResource) NoResource() WithoutResource {
	// Since a domainWithResource is guaranteed to contain a resource, we can assume the slash is there
	slash := strings.IndexRune(string(j), '/')
	return Domain(j[:slash])
}

// Split implements WithResource
func (j domainWithResource) Split() (WithoutResource, Resource) {
	return j.NoResource(), j.Resource()
}

func (j bare) _ForcedToNotHaveResource()   {}
func (j Domain) _ForcedToNotHaveResource() {}

// NoResource implements Any
func (j Domain) NoResource() WithoutResource {
	return j
}

// WithResource implements WithoutResource
func (j Domain) WithResource(r Resource) WithResource {
	return R(j.String() + "/" + string(r))
}

// MaybeWithResource implements WithoutResource
func (j Domain) MaybeWithResource(r Resource) Any {
	return Parse(j.String() + "/" + string(r))
}

// MaybeLocal returns the local part of a JID if it has one, otherwise empty
func MaybeLocal(j Any) Local {
	if jj, ok := j.(WithLocal); ok {
		return jj.Local()
	}
	return Local("")
}

// WithAndWithout will return the JID with the resource, and without the resource
func WithAndWithout(peer Any) (WithResource, WithoutResource) {
	if pwr, ok := peer.(WithResource); ok {
		return pwr, peer.NoResource()
	}
	return nil, peer.NoResource()
}
