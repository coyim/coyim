package muc

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/coyim/coyim/xmpp/jid"
)

type RoomRoster struct {
	sync.RWMutex

	occupants map[string]*Occupant
}

func (r *RoomRoster) occupantList() []*Occupant {
	r.RLock()
	defer r.RUnlock()

	result := []*Occupant{}

	for _, o := range r.occupants {
		result = append(result, o)
	}

	return result
}

func (r *RoomRoster) AllOccupants() []*Occupant {
	result := r.occupantList()
	sort.Sort(ByOccupantNick(result))
	return result
}

func (r *RoomRoster) NoRole() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Role.(*noneRole); ok {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) Visitors() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Role.(*visitorRole); ok {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) Participants() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Role.(*participantRole); ok {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) Moderators() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Role.(*moderatorRole); ok {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) NoAffiliation() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*noneAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) Banned() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*outcastAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) Members() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*memberAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) Admins() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*adminAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) Owners() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*ownerAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) OccupantsByRole() (none, visitors, participants, moderators []*Occupant) {
	return r.NoRole(), r.Visitors(), r.Participants(), r.Moderators()
}

func (r *RoomRoster) OccupantsByAffiliation() (none, banned, members, admins, owners []*Occupant) {
	return r.NoAffiliation(), r.Banned(), r.Members(), r.Admins(), r.Owners()
}

// UpdateNick should be called when receiving an unavailable with status code 303
func (r *RoomRoster) UpdateNick(from jid.WithResource, newNick string) error {
	r.Lock()
	defer r.Unlock()

	base := from.NoResource()
	newFull := base.WithResource(jid.Resource(newNick))

	oc, ok := r.occupants[from.String()]
	if !ok {
		return errors.New("no such occupant known in this room")
	}

	oc.Nick = newNick
	delete(r.occupants, from.String())
	r.occupants[newFull.String()] = oc

	return nil
}

func (r *RoomRoster) UpdatePresence(from jid.WithResource, tp, affiliation, role, show, statusCode, statusMsg string, realJid jid.WithResource) (joined, left bool, err error) {
	r.Lock()
	defer r.Unlock()

	oc, ok := r.occupants[from.String()]

	if tp == "unavailable" {
		if !ok {
			return false, false, errors.New("no such occupant known in this room")
		}
		oc.ChangeRoleToNone()
		delete(r.occupants, from.String())
		return false, true, nil
	}

	if tp != "" {
		return false, false, fmt.Errorf("incorrect presence type sent to room roster: '%s'", tp)
	}

	if !ok {
		oc = &Occupant{}
		err = oc.Update(from, affiliation, role, show, statusMsg, realJid)

		if err != nil {
			return false, false, err
		}

		r.occupants[from.String()] = oc

		return true, false, nil
	}

	err = oc.Update(from, affiliation, role, show, statusMsg, realJid)
	if err != nil {
		return false, false, err
	}
	return false, false, nil
}
