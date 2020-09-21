package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	roomViewRosterStatusIconIndex int = iota
	roomViewRosterNickNameIndex
	roomViewRosterAffiliationIndex
	roomViewRosterRoleIndex
)

type roomViewRoster struct {
	r *muc.RoomRoster

	view  gtki.Box      `gtk-widget:"roomRosterBox"`
	tree  gtki.TreeView `gtk-widget:"room-members-tree"`
	model gtki.ListStore
}

func (v *roomView) newRoomViewRoster() *roomViewRoster {
	r := &roomViewRoster{
		r: v.room.Roster(),
	}

	builder := newBuilder("MUCRoomRoster")
	panicOnDevError(builder.bindObjects(r))

	var err error
	r.model, err = g.gtk.ListStoreNew(pixbufType(), glibi.TYPE_STRING, glibi.TYPE_STRING, glibi.TYPE_STRING)
	if err != nil {
		panic(err)
	}

	r.tree.SetModel(r.model)

	v.subscribe("roster", occupantSelfJoined, r.onUpdateRoster)
	v.subscribe("roster", occupantJoined, r.onUpdateRoster)
	v.subscribe("roster", occupantUpdated, r.onUpdateRoster)
	v.subscribe("roster", occupantLeft, r.onUpdateRoster)

	return r
}

func (v *roomViewRoster) onUpdateRoster(roomViewEventInfo) {
	v.updateRosterModel()
}

func (v *roomViewRoster) updateRosterModel() {
	v.model.Clear()

	for _, o := range v.r.AllOccupants() {
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
	case muc.AffiliationOutcast:
		return i18n.Local("Outcast")
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
