package gui

import (
	"errors"
	"strings"

	"github.com/coyim/gotk3adapter/gdki"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	coyroster "github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	roomViewRosterStatusIconIndex int = iota
	roomViewRosterNicknameIndex
	roomViewRosterAffiliationIndex
	roomViewRosterInfoIndex
	roomViewRosterCanShowInfoIndex
)

type roomViewRoster struct {
	roster *muc.RoomRoster

	areSignalsEnabled bool

	view        gtki.Box      `gtk-widget:"roster-view"`
	rosterPanel gtki.Box      `gtk-widget:"roster-main-panel"`
	tree        gtki.TreeView `gtk-widget:"roster-tree-view"`
	rosterInfo  *roomViewRosterInfo

	model gtki.TreeStore

	log coylog.Logger
}

func (v *roomView) newRoomViewRoster() *roomViewRoster {
	r := &roomViewRoster{
		roster:            v.room.Roster(),
		log:               v.log,
		areSignalsEnabled: true,
	}

	r.initBuilder()
	r.initDefaults()
	r.initSubscribers(v)

	return r
}

func (r *roomViewRoster) initBuilder() {
	builder := newBuilder("MUCRoomRoster")
	builder.ConnectSignals(map[string]interface{}{
		"on_occupant_selected": r.onOccupantSelected,
	})

	panicOnDevError(builder.bindObjects(r))
}

func (r *roomViewRoster) initDefaults() {
	r.model, _ = g.gtk.TreeStoreNew(
		// icon
		pixbufType(),
		// display nickname
		glibi.TYPE_STRING,
		// affiliation
		glibi.TYPE_STRING,
		// info tooltip
		glibi.TYPE_STRING,
		// can show information
		glibi.TYPE_BOOLEAN,
	)

	r.tree.SetModel(r.model)
	r.draw()
}

func (r *roomViewRoster) initSubscribers(v *roomView) {
	v.subscribe("roster", func(ev roomViewEvent) {
		switch ev.(type) {
		case occupantSelfJoinedEvent:
			r.onUpdateRoster()
		case occupantJoinedEvent:
			r.onUpdateRoster()
		case occupantUpdatedEvent:
			r.onUpdateRoster()
		case occupantLeftEvent:
			r.onUpdateRoster()
		case selfOccupantRemovedEvent:
			r.onUpdateRoster()
		case occupantRemovedEvent:
			r.onUpdateRoster()
		case roomDestroyedEvent:
			r.onRoomDestroy()
		}
	})
}

func (r *roomViewRoster) onOccupantSelected(_ gtki.TreeView, path gtki.TreePath) {
	if !r.areSignalsEnabled {
		return
	}

	nickname, err := r.getNicknameFromTreeModel(path)
	if err != nil {
		r.log.Warn("Nickname not found")
		return
	}

	occupant, ok := r.roster.GetOccupant(nickname)
	if !ok {
		r.log.WithField("nickname", nickname).Warn("Occupant was not found")
		return
	}

	r.showOccupantInfo(occupant)
}

func (r *roomViewRoster) addInfoPanel() {
	r.rosterInfo = r.newRoomViewRosterInfo()
}

func (r *roomViewRoster) showOccupantInfo(occupant *muc.Occupant) {
	r.rosterPanel.Hide()

	r.addInfoPanel()
	r.view.Add(r.rosterInfo.rosterInfoBox)
	r.rosterInfo.displayOccupantInfoPanel(occupant, func() {
		doInUIThread(func() {
			r.rosterPanel.Show()
		})
	})
}

func (r *roomViewRoster) getNicknameFromTreeModel(path gtki.TreePath) (string, error) {
	iter, err := r.model.GetIter(path)
	if err != nil {
		return "", err
	}

	nickname, err := r.getOccupantInfoFromIter(iter)
	if err != nil {
		return "", err
	}

	return nickname, nil
}

func (r *roomViewRoster) getOccupantInfoFromIter(iter gtki.TreeIter) (string, error) {
	canShowInfoValue, err := r.getIterValue(iter, roomViewRosterCanShowInfoIndex)
	if err != nil {
		return "", err
	}

	if !canShowInfoValue.(bool) {
		return "", errors.New("The node selected is not an occupant")
	}

	nickname, err := r.getIterValue(iter, roomViewRosterNicknameIndex)
	if err != nil {
		return "", err
	}

	return nickname.(string), nil
}

func (r *roomViewRoster) getIterValue(iter gtki.TreeIter, index int) (interface{}, error) {
	iterValue, e1 := r.model.GetValue(iter, index)
	iterRef, e2 := iterValue.GoValue()
	if e1 != nil || e2 != nil {
		return 0, errors.New("error trying to get iter value")
	}

	return iterRef, nil
}

func (r *roomViewRoster) onUpdateRoster() {
	doInUIThread(r.redraw)
}

func (r *roomViewRoster) onRoomDestroy() {
	r.areSignalsEnabled = false
}
func (r *roomViewRoster) draw() {
	noneRoles, visitors, participants, moderators := r.roster.OccupantsByRole()

	r.drawOccupantsByRole(data.RoleModerator, moderators)
	r.drawOccupantsByRole(data.RoleParticipant, participants)
	r.drawOccupantsByRole(data.RoleVisitor, visitors)
	r.drawOccupantsByRole(data.RoleNone, noneRoles)

	r.tree.ExpandAll()
}

func (r *roomViewRoster) redraw() {
	r.model.Clear()
	r.draw()
}

func (r *roomViewRoster) drawOccupantsByRole(role string, occupants []*muc.Occupant) {
	if isOccupantListEmpty(occupants) {
		return
	}

	roleHeader := rolePluralName(role)
	roleHeader = i18n.Localf("%s (%v)", roleHeader, len(occupants))

	iter := r.model.Append(nil)
	r.model.SetValue(iter, roomViewRosterNicknameIndex, roleHeader)
	r.model.SetValue(iter, roomViewRosterCanShowInfoIndex, false)

	for _, o := range occupants {
		r.addOccupantToRoster(o, iter)
	}
}

func (r *roomViewRoster) addOccupantToRoster(o *muc.Occupant, parentIter gtki.TreeIter) {
	iter := r.model.Append(parentIter)

	r.model.SetValue(iter, roomViewRosterStatusIconIndex, getOccupantIconForStatus(o.Status))
	r.model.SetValue(iter, roomViewRosterNicknameIndex, o.Nickname)
	r.model.SetValue(iter, roomViewRosterAffiliationIndex, affiliationDisplayName(o.Affiliation))
	r.model.SetValue(iter, roomViewRosterInfoIndex, occupantDisplayTooltip(o))
	r.model.SetValue(iter, roomViewRosterCanShowInfoIndex, true)
}

func getOccupantIconForStatus(s *coyroster.Status) gdki.Pixbuf {
	icon := getOccupantIconNameForStatus(s.Status)
	return getMUCIconPixbuf(icon)
}

func getOccupantIconNameForStatus(status string) string {
	switch status {
	case "unavailable":
		return "occupant-offline"
	case "away":
		return "occupant-away"
	case "dnd":
		return "occupant-busy"
	case "xa":
		return "occupant-extended-away"
	case "chat":
		return "occupant-chat"
	default:
		return "occupant-online"
	}
}

func affiliationDisplayName(a data.Affiliation) string {
	switch a.Name() {
	case data.AffiliationAdmin:
		return i18n.Local("Admin")
	case data.AffiliationOwner:
		return i18n.Local("Owner")
	case data.AffiliationOutcast:
		return i18n.Local("Outcast")
	default: // Member or other values get the default treatment
		return ""
	}
}

func rolePluralName(role string) string {
	switch role {
	case data.RoleNone:
		return i18n.Local("None")
	case data.RoleParticipant:
		return i18n.Local("Participants")
	case data.RoleVisitor:
		return i18n.Local("Visitors")
	case data.RoleModerator:
		return i18n.Local("Moderators")
	default:
		// This should not really be possible, but it is necessary
		// because golang can't prove it
		return ""
	}
}

func roleDisplayName(role data.Role) string {
	switch role.Name() {
	case data.RoleNone:
		return i18n.Local("None")
	case data.RoleParticipant:
		return i18n.Local("Participant")
	case data.RoleVisitor:
		return i18n.Local("Visitor")
	case data.RoleModerator:
		return i18n.Local("Moderator")
	default:
		// This should not really be possible, but it is necessary
		// because golang can't prove it
		return ""
	}
}

func statusDisplayMessage(s *coyroster.Status) string {
	return s.StatusMsg
}

func statusDisplayName(s *coyroster.Status) string {
	return showForDisplay(s.Status, false)
}

func occupantDisplayTooltip(o *muc.Occupant) string {
	ms := []string{
		o.Nickname,
		statusDisplayName(o.Status),
	}

	m := statusDisplayMessage(o.Status)
	if m != "" {
		ms = append(ms, m)
	}

	return strings.Join(ms, "\n")
}

func isOccupantListEmpty(o []*muc.Occupant) bool {
	return len(o) == 0
}
