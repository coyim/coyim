package gui

import (
	"runtime"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	roomViewRosterStatusIconIndex int = iota
	roomViewRosterNickNameIndex
	roomViewRosterAffiliationIndex
	roomViewRosterRoleIndex
)

type roomViewRoster struct {
	view  gtki.Box       `gtk-widget:"roomRosterBox"`
	model gtki.ListStore `gtk-widget:"room-members-model"`
	tree  gtki.TreeView  `gtk-widget:"room-members-tree"`
}

func newRoomViewRoster() *roomViewRoster {
	r := &roomViewRoster{}

	builder := newBuilder("MUCRoomRoster")
	panicOnDevError(builder.bindObjects(r))

	// r.model needs to be kept beyond the lifespan of the builder.
	r.model.Ref()
	runtime.SetFinalizer(r, func(ros interface{}) {
		ros.(*roster).model.Unref()
		ros.(*roster).model = nil
	})

	return r
}

func (v *roomViewRoster) updateRoomRoster(occupants []*muc.Occupant) {
	v.model.Clear()

	for _, o := range occupants {
		v.addOccupantToRoster(o)
	}

	v.tree.ExpandAll()
}

func (v *roomViewRoster) addOccupantToRoster(o *muc.Occupant) {
	iter := v.model.Append()

	_ = v.model.SetValue(iter, roomViewRosterStatusIconIndex, v.getOccupantIcon().GetPixbuf())
	_ = v.model.SetValue(iter, roomViewRosterNickNameIndex, o.Nick)
	_ = v.model.SetValue(iter, roomViewRosterAffiliationIndex, v.affiliationDisplayName(o.Affiliation))
	_ = v.model.SetValue(iter, roomViewRosterRoleIndex, v.roleDisplayName(o.Role))
}

func (v *roomViewRoster) getOccupantIcon() Icon {
	return statusIcons["occupant"]
}

func (v *roomViewRoster) affiliationDisplayName(a muc.Affiliation) string {
	switch a.Name() {
	case muc.AffiliationAdmin:
		return i18n.Local("Admin")
	case muc.AffiliationOwner:
		return i18n.Local("Owner")
	default:
		return ""
	}
}

func (v *roomViewRoster) roleDisplayName(r muc.Role) string {
	switch r.Name() {
	case muc.RoleNone:
		return i18n.Local("None")
	case muc.RoleParticipant:
		return i18n.Local("Participant")
	case muc.RoleVisitor:
		return i18n.Local("Visitor")
	case muc.RoleModerator:
		return i18n.Local("Moderator")
	default:
		return ""
	}
}
