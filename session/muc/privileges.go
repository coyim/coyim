package muc

import (
	"sync"

	"github.com/coyim/coyim/session/muc/data"
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

func definedPrivilegesForGroup(groupName string) map[string]*privileges {
	privilegesInitOnce.Do(func() {
		privilegesSingleton = map[string]map[string]*privileges{
			"roles":        definedPrivilegesForRoles(),
			"affiliations": definedPrivilegesForAffiliations(),
		}
	})
	return privilegesSingleton[groupName]
}

func definedPrivilegesForGroupItem(groupName, groupItem string) *privileges {
	privilegesGroup := definedPrivilegesForGroup(groupName)
	return privilegesGroup[groupItem]
}

func definedPrivilegesForRole(r data.Role) *privileges {
	return definedPrivilegesForGroupItem("roles", r.Name())
}

func definedPrivilegesForAffiliation(a data.Affiliation) *privileges {
	return definedPrivilegesForGroupItem("affiliations", a.Name())
}
