package muc

import (
	"sync"

	"github.com/golang-collections/collections/set"
)

type privilege int

var (
	privilegesInitOnce  sync.Once
	privilegesSingleton map[string]map[string]*privileges
)

type privileges struct {
	list *set.Set
}

func newPrivileges(l ...privilege) *privileges {
	p := &privileges{
		list: set.New(),
	}

	for _, px := range l {
		p.list.Insert(px)
	}

	return p
}

func (p *privileges) can(privilege privilege) bool {
	return p.list.Has(privilege)
}
