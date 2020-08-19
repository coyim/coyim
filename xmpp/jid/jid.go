package jid

import (
	"fmt"
	"strings"
)

// Local represents the local part of a JID
type Local struct {
	v string
}

// Domain represents the domain part of a JID
type Domain struct {
	v string
}

// Resource represents the resource part of a JID
type Resource struct {
	v string
}

// NewLocal returns a new local if possible
func NewLocal(s string) Local {
	if !ValidLocal(s) {
		return Local{""}
	}

	return Local{s}
}

// NewDomain returns a new domain if possible
func NewDomain(s string) Domain {
	if !ValidDomain(s) {
		return Domain{""}
	}

	return Domain{s}
}

// NewResource returns a new resource if possible
func NewResource(s string) Resource {
	if !ValidResource(s) {
		return Resource{""}
	}

	return Resource{s}
}

type bare struct {
	l Local
	d Domain
}

type full struct {
	l Local
	d Domain
	r Resource
}

type domainWithResource struct {
	d Domain
	r Resource
}

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
	// Valid returns true if this is a valid JID
	Valid() bool
}

// NewBare accept a domain and a local and creates a valid resource
func NewBare(local Local, domain Domain) Bare {
	return domain.AddLocal(local)
}

// NewFull creates a full JID from the different parts of a JID
func NewFull(local Local, domain Domain, resource Resource) Full {
	return domain.AddLocal(local).WithResource(resource).(Full)
}

// WithResource represents any valid JID that has a resource part
type WithResource interface {
	Any
	// Resource will return the resource
	Resource() Resource
	// Split will return the JID split into the part without resource and the part with resource
	Split() (WithoutResource, Resource)
	// Bare returns a bare jid from the actual jid
	Bare() Bare
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

// ParseBare returns a bare JID. It will fail if the given string isn't at least a bare
func ParseBare(s string) Bare {
	return NR(s).(Bare)
}

// TryParseBare returns a bare JID if it can.
func TryParseBare(s string) (Bare, bool) {
	res, ok := NR(s).(Bare)
	return res, ok
}

// ParseFull returns a full JID. It will fail if the given string isn't at least a full
func ParseFull(s string) Full {
	return R(s).(Full)
}

// TryParseFull returns a full JID if it can.
func TryParseFull(s string) (Full, bool) {
	res, ok := R(s).(Full)
	return res, ok
}

// ParseDomain returns a domain part of a JID. It will fail if the given string isn't at least a domain
// This will parse the full string as a JID and _extract_ the domain part, This is in comparison to
// NewDomain that will try to create a new Domain object from the given string
func ParseDomain(s string) Domain {
	return NR(s).Host()
}

// Parse will parse the given string and return the most specific JID type that matches it
// In general, it is a good idea to check that the returned result is valid before using it by calling Valid()
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
			return full{NewLocal(local), NewDomain(left), NewResource(resource)}
		}
		return bare{NewLocal(local), NewDomain(left)}
	}

	if resource != "" {
		return domainWithResource{NewDomain(left), NewResource(resource)}
	}

	return NewDomain(left)
}

// Host implements Any
func (j Domain) Host() Domain {
	return j
}

// String implements Any
func (j Domain) String() string {
	return j.v
}

// PotentialResource implements Any
func (j Domain) PotentialResource() Resource {
	return Resource{""}
}

// PotentialSplit implements Any
func (j Domain) PotentialSplit() (WithoutResource, Resource) {
	return j, j.PotentialResource()
}

// Host implements Any
func (j bare) Host() Domain {
	return j.d
}

// String implements Any
func (j bare) String() string {
	return fmt.Sprintf("%s@%s", j.l, j.d)
}

// PotentialResource implements Any
func (j bare) PotentialResource() Resource {
	return Resource{""}
}

// PotentialSplit implements Any
func (j bare) PotentialSplit() (WithoutResource, Resource) {
	return j, j.PotentialResource()
}

// Local implements Bare
func (j bare) Local() Local {
	return j.l
}

// NoResource implements Any
func (j bare) NoResource() WithoutResource {
	return j
}

// WithResource implements Any
func (j bare) WithResource(r Resource) WithResource {
	return R(j.String() + "/" + r.v)
}

// MaybeWithResource implements Any
func (j bare) MaybeWithResource(r Resource) Any {
	return Parse(j.String() + "/" + r.v)
}

// Host implements Any
func (j full) Host() Domain {
	return j.NoResource().Host()
}

// String implements Any
func (j full) String() string {
	return fmt.Sprintf("%s@%s/%s", j.l, j.d, j.r)
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
	return bare{j.l, j.d}
}

// Resource implements WithResource
func (j full) Resource() Resource {
	return j.r
}

// Split implements WithResource
func (j full) Split() (WithoutResource, Resource) {
	return j.NoResource(), j.Resource()
}

func (j full) Bare() Bare {
	return j.NoResource().(bare)
}

// Host implements Any
func (j domainWithResource) Host() Domain {
	return NewDomain(j.NoResource().String())
}

// String implements Any
func (j domainWithResource) String() string {
	return fmt.Sprintf("%s/%s", j.d, j.r)
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
	return j.r
}

// NoResource implements WithResource
func (j domainWithResource) NoResource() WithoutResource {
	return j.d
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
	return R(fmt.Sprintf("%s/%s", j, r))
}

// MaybeWithResource implements WithoutResource
func (j Domain) MaybeWithResource(r Resource) Any {
	return Parse(fmt.Sprintf("%s/%s", j, r))
}

// AddLocal returns a Bare, combining this domain with a Local
func (j Domain) AddLocal(l Local) Bare {
	return ParseBare(fmt.Sprintf("%s@%s", l, j))
}

// MaybeLocal returns the local part of a JID if it has one, otherwise empty
func MaybeLocal(j Any) Local {
	if jj, ok := j.(WithLocal); ok {
		return jj.Local()
	}
	return Local{""}
}

// WithAndWithout will return the JID with the resource, and without the resource
func WithAndWithout(peer Any) (WithResource, WithoutResource) {
	if pwr, ok := peer.(WithResource); ok {
		return pwr, peer.NoResource()
	}
	return nil, peer.NoResource()
}

// String implements Local
func (j Local) String() string {
	return j.v
}

// String implements Resource
func (j Resource) String() string {
	return j.v
}

// Valid returns true if this object is valid
func (j Local) Valid() bool {
	return j.v != ""
}

// Valid returns true if this object is valid
func (j Domain) Valid() bool {
	return j.v != ""
}

// Valid returns true if this object is valid
func (j Resource) Valid() bool {
	return j.v != ""
}

// Valid returns true if this object is valid
func (j bare) Valid() bool {
	return j.l.Valid() && j.d.Valid()
}

// Valid returns true if this object is valid
func (j full) Valid() bool {
	return j.l.Valid() && j.d.Valid() && j.r.Valid()
}

// Valid returns true if this object is valid
func (j domainWithResource) Valid() bool {
	return j.d.Valid() && j.r.Valid()
}
