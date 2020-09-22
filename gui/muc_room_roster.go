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
	roster *muc.RoomRoster

	view  gtki.Box      `gtk-widget:"roomRosterBox"`
	tree  gtki.TreeView `gtk-widget:"room-members-tree"`
	model gtki.ListStore
}

func (v *roomView) newRoomViewRoster() *roomViewRoster {
	r := &roomViewRoster{
		roster: v.room.Roster(),
	}

	r.initBuilder()
	r.initDefaults()
	r.initSubscribers(v)
	return r
}

func (r *roomViewRoster) initBuilder() {
	builder := newBuilder("MUCRoomRoster")
	panicOnDevError(builder.bindObjects(r))
}

func (r *roomViewRoster) initDefaults() {
	var err error
	r.model, err = g.gtk.ListStoreNew(pixbufType(), glibi.TYPE_STRING, glibi.TYPE_STRING, glibi.TYPE_STRING)
	if err != nil {
		panic(err)
	}

	r.tree.SetModel(r.model)
	r.draw()
}

func (r *roomViewRoster) initSubscribers(v *roomView) {
	v.subscribe("roster", occupantSelfJoined, r.onUpdateRoster)
	v.subscribe("roster", occupantJoined, r.onUpdateRoster)
	v.subscribe("roster", occupantUpdated, r.onUpdateRoster)
	v.subscribe("roster", occupantLeft, r.onUpdateRoster)
}

func (r *roomViewRoster) onUpdateRoster(roomViewEventInfo) {
	r.redraw()
}

func (r *roomViewRoster) draw() {
	for _, o := range r.roster.AllOccupants() {
		r.addOccupantToRoster(o)
	}

	r.tree.ExpandAll()
}

func (r *roomViewRoster) redraw() {
	r.model.Clear()
	r.draw()
}

func (r *roomViewRoster) addOccupantToRoster(o *muc.Occupant) {
	iter := r.model.Append()

	_ = r.model.SetValue(iter, roomViewRosterStatusIconIndex, r.getOccupantIcon().GetPixbuf())
	_ = r.model.SetValue(iter, roomViewRosterNickNameIndex, o.Nick)
	_ = r.model.SetValue(iter, roomViewRosterAffiliationIndex, r.affiliationDisplayName(o.Affiliation))
	_ = r.model.SetValue(iter, roomViewRosterRoleIndex, r.roleDisplayName(o.Role))
}

func (r *roomViewRoster) getOccupantIcon() Icon {
	return statusIcons["occupant"]
}

func (r *roomViewRoster) affiliationDisplayName(a muc.Affiliation) string {
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

func (r *roomViewRoster) roleDisplayName(role muc.Role) string {
	switch role.Name() {
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
