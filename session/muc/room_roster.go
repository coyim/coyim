package muc

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/muc/data"

	"github.com/coyim/coyim/xmpp/jid"
)

// OccupantPresenceInfo containts information for an occupant presence
type OccupantPresenceInfo struct {
	Nickname      string
	RealJid       jid.Full
	Affiliation   data.Affiliation
	Role          data.Role
	Status        string
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
	none, _, _, _ := r.OccupantsByRole()
	return none
}

// Visitors returns all occupants that have the visitor role in a room, sorted by nickname
func (r *RoomRoster) Visitors() []*Occupant {
	_, visitors, _, _ := r.OccupantsByRole()
	return visitors
}

// Participants returns all occupants that have the participant role in a room, sorted by nickname
func (r *RoomRoster) Participants() []*Occupant {
	_, _, participants, _ := r.OccupantsByRole()
	return participants
}

// Moderators returns all occupants that have the moderator role in a room, sorted by nickname
func (r *RoomRoster) Moderators() []*Occupant {
	_, _, _, moderators := r.OccupantsByRole()
	return moderators
}

// NoAffiliation returns all occupants that have no affiliation in a room, sorted by nickname
func (r *RoomRoster) NoAffiliation() []*Occupant {
	none, _, _, _, _ := r.OccupantsByAffiliation()
	return none
}

// Banned returns all occupants that are banned in a room, sorted by nickname. This should likely not return anything.
func (r *RoomRoster) Banned() []*Occupant {
	_, banned, _, _, _ := r.OccupantsByAffiliation()
	return banned
}

// Members returns all occupants that are members in a room, sorted by nickname.
func (r *RoomRoster) Members() []*Occupant {
	_, _, members, _, _ := r.OccupantsByAffiliation()
	return members
}

// Admins returns all occupants that are administrators in a room, sorted by nickname.
func (r *RoomRoster) Admins() []*Occupant {
	_, _, _, admins, _ := r.OccupantsByAffiliation()
	return admins
}

// Owners returns all occupants that are owners in a room, sorted by nickname.
func (r *RoomRoster) Owners() []*Occupant {
	_, _, _, _, owners := r.OccupantsByAffiliation()
	return owners
}

// OccupantsByRole returns all occupants, divided into the different roles
func (r *RoomRoster) OccupantsByRole() (none, visitors, participants, moderators []*Occupant) {
	for _, o := range r.AllOccupants() {
		switch o.Role.(type) {
		case *data.NoneRole:
			none = append(none, o)
		case *data.VisitorRole:
			visitors = append(visitors, o)
		case *data.ParticipantRole:
			participants = append(participants, o)
		case *data.ModeratorRole:
			moderators = append(moderators, o)
		}
	}

	return
}

// OccupantsByAffiliation returns all occupants, divided into the different affiliations
func (r *RoomRoster) OccupantsByAffiliation() (none, banned, members, admins, owners []*Occupant) {
	for _, o := range r.AllOccupants() {
		switch o.Affiliation.(type) {
		case *data.NoneAffiliation:
			none = append(none, o)
		case *data.OutcastAffiliation:
			banned = append(banned, o)
		case *data.MemberAffiliation:
			members = append(members, o)
		case *data.AdminAffiliation:
			admins = append(admins, o)
		case *data.OwnerAffiliation:
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
			Status:    op.Status,
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

	o.Update(op.Nickname, op.Affiliation, op.Role, op.Status, op.StatusMessage, op.RealJid)
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
