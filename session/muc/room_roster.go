package muc

import (
	"errors"
	"fmt"
	"reflect"
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

func (r *RoomRoster) byRole(role interface{}) []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if areTheSameType(role, o.Role) {
			result = append(result, o)
		}
	}

	return result
}

func (r *RoomRoster) byAffiliation(affiliation interface{}) []*Occupant {
	result := []*Occupant{}

	for _, o := range r.AllOccupants() {
		if areTheSameType(affiliation, o.Affiliation) {
			result = append(result, o)
		}
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
	return r.byRole(&noneRole{})
}

// Visitors returns all occupants that have the visitor role in a room, sorted by nickname
func (r *RoomRoster) Visitors() []*Occupant {
	return r.byRole(&visitorRole{})
}

// Participants returns all occupants that have the participant role in a room, sorted by nickname
func (r *RoomRoster) Participants() []*Occupant {
	return r.byRole(&participantRole{})
}

// Moderators returns all occupants that have the moderator role in a room, sorted by nickname
func (r *RoomRoster) Moderators() []*Occupant {
	return r.byRole(&moderatorRole{})
}

// NoAffiliation returns all occupants that have no affiliation in a room, sorted by nickname
func (r *RoomRoster) NoAffiliation() []*Occupant {
	return r.byAffiliation(&noneAffiliation{})
}

// Banned returns all occupants that are banned in a room, sorted by nickname. This should likely not return anything.
func (r *RoomRoster) Banned() []*Occupant {
	return r.byAffiliation(&outcastAffiliation{})
}

// Members returns all occupants that are members in a room, sorted by nickname.
func (r *RoomRoster) Members() []*Occupant {
	return r.byAffiliation(&memberAffiliation{})
}

// Admins returns all occupants that are administrators in a room, sorted by nickname.
func (r *RoomRoster) Admins() []*Occupant {
	return r.byAffiliation(&adminAffiliation{})
}

// Owners returns all occupants that are owners in a room, sorted by nickname.
func (r *RoomRoster) Owners() []*Occupant {
	return r.byAffiliation(&ownerAffiliation{})
}

// OccupantsByRole returns all occupants, divided into the different roles
func (r *RoomRoster) OccupantsByRole() (none, visitors, participants, moderators []*Occupant) {
	for _, o := range r.AllOccupants() {
		switch o.Role.(type) {
		case *noneRole:
			none = append(none, o)
		case *visitorRole:
			visitors = append(visitors, o)
		case *participantRole:
			participants = append(participants, o)
		case *moderatorRole:
			moderators = append(moderators, o)
		}
	}

	return
}

// OccupantsByAffiliation returns all occupants, divided into the different affiliations
func (r *RoomRoster) OccupantsByAffiliation() (none, banned, members, admins, owners []*Occupant) {
	for _, o := range r.AllOccupants() {
		switch o.Affiliation.(type) {
		case *noneAffiliation:
			none = append(none, o)
		case *outcastAffiliation:
			banned = append(banned, o)
		case *memberAffiliation:
			members = append(members, o)
		case *adminAffiliation:
			admins = append(admins, o)
		case *ownerAffiliation:
			owners = append(owners, o)
		}
	}

	return
}

// UpdateNickname should be called when receiving an unavailable with status code 303
// The new nickname should be given without the room name
func (r *RoomRoster) UpdateNickname(nickname, newNickname string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	oc, ok := r.occupants[nickname]
	if !ok {
		return errors.New("no such occupant known in this room")
	}

	oc.Nickname = newNickname
	delete(r.occupants, nickname)
	r.occupants[newNickname] = oc

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
		RealJid:     op.RealJid,
		Nickname:    op.Nickname,
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
		r.occupants[o.Nickname] = o
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

func areTheSameType(v1 interface{}, v2 interface{}) bool {
	return reflect.TypeOf(v1) == reflect.TypeOf(v2)
}
