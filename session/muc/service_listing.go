package muc

import (
	"github.com/coyim/coyim/xmpp/jid"
)

// ServiceListing contains the information about a service for listing it
type ServiceListing struct {
	Jid  jid.Any
	Name string
}

// NewServiceListing creates and returns a new service listing
func NewServiceListing(service jid.Any, name string) *ServiceListing {
	sl := &ServiceListing{
		Jid:  service,
		Name: name,
	}
	return sl
}
