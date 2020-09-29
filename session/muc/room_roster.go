package muc

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/coyim/coyim/roster"

	"github.com/coyim/coyim/xmpp/jid"
)

// OccupantPresenceInfo containts information for an occupant presence
type OccupantPresenceInfo struct {
	Nickname      string
	RealJid       jid.Full
	Affiliation   Affiliation
	Role          Role
	Show          string
	StatusCode    string
	StatusMessage string
}

// RoomRoster contains information about all the occupants in a room
type RoomRoster struct {
	lock sync.RWMutex

	occupants map[string]*Occupant
}

// newRoomRoster returns a newly created room roster
func newRoomRoster() *RoomRoster {
	return &RoomRoster{
		occupants: make(map[string]*Occupant),
	}
}

func (r *RoomRoster) occupantList() []*Occupant {
	r.lock.RLock()
	defer r.lock.RUnlock()

	result := []*Occupant{}

	for _, o := range r.occupants {
		result = append(result, o)
	}

	return result
}

// AllOccupants returns a list of all occupants in the room, sorted by nickname
func (r *RoomRoster) AllOccupants() []*Occupant {
	result := r.occupantList()
	sort.Sort(ByOccupantNick(result))
	return result
}

// NoRole returns all occupants that have no role in a room, sorted by nickname
func (r *RoomRoster) NoRole() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Role.(*noneRole); ok {
			result = append(result, o)
		}
	}

	return result
}

// Visitors returns all occupants that have the visitor role in a room, sorted by nickname
func (r *RoomRoster) Visitors() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Role.(*visitorRole); ok {
			result = append(result, o)
		}
	}

	return result
}

// Participants returns all occupants that have the participant role in a room, sorted by nickname
func (r *RoomRoster) Participants() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Role.(*participantRole); ok {
			result = append(result, o)
		}
	}

	return result
}

// Moderators returns all occupants that have the moderator role in a room, sorted by nickname
func (r *RoomRoster) Moderators() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Role.(*moderatorRole); ok {
			result = append(result, o)
		}
	}

	return result
}

// NoAffiliation returns all occupants that have no affiliation in a room, sorted by nickname
func (r *RoomRoster) NoAffiliation() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*noneAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

// Banned returns all occupants that are banned in a room, sorted by nickname. This should likely not return anything.
func (r *RoomRoster) Banned() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*outcastAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

// Members returns all occupants that are members in a room, sorted by nickname.
func (r *RoomRoster) Members() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*memberAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

// Admins returns all occupants that are administrators in a room, sorted by nickname.
func (r *RoomRoster) Admins() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*adminAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

// Owners returns all occupants that are owners in a room, sorted by nickname.
func (r *RoomRoster) Owners() []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if _, ok := o.Affiliation.(*ownerAffiliation); ok {
			result = append(result, o)
		}
	}

	return result
}

// OccupantsByRole returns all occupants, divided into the different roles
func (r *RoomRoster) OccupantsByRole() (none, visitors, participants, moderators []*Occupant) {
	// TODO: This should be made thread safe
	return r.NoRole(), r.Visitors(), r.Participants(), r.Moderators()
}

// OccupantsByAffiliation returns all occupants, divided into the different affiliations
func (r *RoomRoster) OccupantsByAffiliation() (none, banned, members, admins, owners []*Occupant) {
	// TODO: This should be made thread safe
	return r.NoAffiliation(), r.Banned(), r.Members(), r.Admins(), r.Owners()
}

// UpdateNick should be called when receiving an unavailable with status code 303
// The new nickname should be given without the room name
func (r *RoomRoster) UpdateNick(from jid.WithResource, newNick string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// TODO: from should not be WithResource, it should be jid.Full

	base := from.NoResource()
	newFull := base.WithResource(jid.NewResource(newNick))

	// TODO: Occupants is only stored with nicknames, so the below is not correct anymore
	oc, ok := r.occupants[from.String()]
	if !ok {
		return errors.New("no such occupant known in this room")
	}

	oc.Nick = newNick
	delete(r.occupants, from.String())
	r.occupants[newFull.String()] = oc

	return nil
}

// UpdatePresence should be called when receiving a regular presence update with no type, or with unavailable as type. It will return
// indications on whether the presence update means the person joined the room, or left the room.
// Notice that updating of nick names is done separately and should not be done by calling this method.
func (r *RoomRoster) UpdatePresence(op *OccupantPresenceInfo, tp string) (joined, left bool, err error) {
	if len(op.Nickname) == 0 {
		return false, false, errors.New("nickname was not provided")
	}

	switch tp {
	case "unavailable":
		err := r.RemoveOccupant(op.Nickname)
		return false, err == nil, err
	case "":
		updated := r.UpdateOrAddOccupant(op)
		return !updated, false, err
	default:
		return false, false, fmt.Errorf("incorrect presence type sent to room roster: '%s'", tp)
	}
}

func (r *RoomRoster) newOccupantFromPresenceInfo(op *OccupantPresenceInfo) *Occupant {
	return &Occupant{
		Affiliation: op.Affiliation,
		Jid:         op.RealJid,
		Nick:        op.Nickname,
		Role:        op.Role,
		Status: &roster.Status{
			Status:    op.StatusCode,
			StatusMsg: op.StatusMessage,
		},
	}
}

// UpdateOrAddOccupant return true if the occupant was updated or false if that was added
func (r *RoomRoster) UpdateOrAddOccupant(op *OccupantPresenceInfo) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	o, ok := r.occupants[op.Nickname]
	if !ok {
		o = r.newOccupantFromPresenceInfo(op)
		r.occupants[o.Nick] = o
		return false
	}

	o.Update(op.Nickname, op.Affiliation, op.Role, op.StatusCode, op.StatusMessage, op.RealJid)
	return true
}

// RemoveOccupant delete an occupant if that exists
func (r *RoomRoster) RemoveOccupant(nickname string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	o, ok := r.occupants[nickname]
	if !ok {
		return errors.New("no such occupant known in this room")
	}

	o.ChangeRoleToNone()
	o.UpdateStatus("unavailable", "Occupant left the room")
	delete(r.occupants, nickname)

	return nil
}

// GetOccupant return an occupant if this exist in the roster, otherwise return nil and false
func (r *RoomRoster) GetOccupant(nickname string) (*Occupant, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	o, ok := r.occupants[nickname]
	return o, ok
}
