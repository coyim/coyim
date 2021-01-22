package gui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coyim/gotk3adapter/gdki"
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	coyroster "github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	roomViewRosterStatusIconIndex int = iota
	roomViewRosterNicknameIndex
	roomViewRosterAffiliationIndex
	roomViewRosterInfoIndex
)

type roomViewRoster struct {
	u        *gtkUI
	roomView *roomView

	roster  *muc.RoomRoster
	account *account
	roomID  jid.Bare

	view        gtki.Box      `gtk-widget:"roster-view"`
	rosterPanel gtki.Box      `gtk-widget:"roster-main-panel"`
	tree        gtki.TreeView `gtk-widget:"roster-tree-view"`

	model      gtki.TreeStore
	rosterInfo *roomViewRosterInfo

	log coylog.Logger
}

func (v *roomView) newRoomViewRoster() *roomViewRoster {
	r := &roomViewRoster{
		u:        v.u,
		roomView: v,
		roster:   v.room.Roster(),
		account:  v.account,
		roomID:   v.roomID(),
		log:      v.log,
	}

	r.initBuilder()
	r.initRosterInfo()
	r.initDefaults()
	r.initSubscribers()

	return r
}

func (r *roomViewRoster) initBuilder() {
	builder := newBuilder("MUCRoomRoster")
	builder.ConnectSignals(map[string]interface{}{
		"on_occupant_selected": r.onOccupantSelected,
	})

	panicOnDevError(builder.bindObjects(r))
}

func (r *roomViewRoster) initRosterInfo() {
	r.rosterInfo = r.newRoomViewRosterInfo()
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
	)

	r.tree.SetModel(r.model)
	r.draw()
}

func (r *roomViewRoster) initSubscribers() {
	r.roomView.subscribe("roster", func(ev roomViewEvent) {
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
		}
	})
}

func (r *roomViewRoster) onOccupantSelected(_ gtki.TreeView, path gtki.TreePath) {
	nickname, err := r.getNicknameFromTreeModel(path)
	if err != nil {
		r.log.Warn("Nickname not found")
		return
	}

	o, ok := r.roster.GetOccupant(nickname)
	if !ok {
		r.log.WithField("nickname", nickname).Debug("Occupant was not found")
		return
	}

	r.showOccupantInfo(o)
}

// updateOccupantAffiliation MUST NOT be called from the UI thread
func (r *roomViewRoster) updateOccupantAffiliation(o *muc.Occupant, previousAffiliation data.Affiliation, reason string) {
	r.log.WithFields(log.Fields{
		"where":       "occupantAffiliationUpdate",
		"occupant":    fmt.Sprintf("%s", o.RealJid),
		"affiliation": o.Affiliation.Name(),
	}).Info("The occupant affiliation has been updated")

	r.roomView.tryUpdateOccupantAffiliation(o, previousAffiliation, reason)
}

// showOccupantInfo MUST be called from the UI thread
func (r *roomViewRoster) showOccupantInfo(o *muc.Occupant) {
	r.rosterInfo.showOccupantInfo(o)
	r.showRosterInfoPanel()
}

// showRosterInfoPanel MUST be called from the UI thread
func (r *roomViewRoster) showRosterInfoPanel() {
	r.rosterPanel.Hide()
	r.view.Add(r.rosterInfo.widget())
}

// hideRosterInfoPanel MUST be called from the UI thread
func (r *roomViewRoster) hideRosterInfoPanel() {
	r.view.Remove(r.rosterInfo.widget())
	r.rosterPanel.Show()
}

func (r *roomViewRoster) getNicknameFromTreeModel(path gtki.TreePath) (string, error) {
	iter, err := r.model.GetIter(path)
	if err != nil {
		r.log.WithError(err).Error("Couldn't activate the selected item")
		return "", err
	}

	iterValue, e1 := r.model.GetValue(iter, roomViewRosterNicknameIndex)
	if e1 != nil {
		return "", errors.New("error trying to get iter value")
	}

	return iterValue.GetString()
}

func (r *roomViewRoster) onUpdateRoster() {
	doInUIThread(r.redraw)
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
	_ = r.model.SetValue(iter, roomViewRosterNicknameIndex, roleHeader)

	for _, o := range occupants {
		r.addOccupantToRoster(o, iter)
	}
}

func (r *roomViewRoster) addOccupantToRoster(o *muc.Occupant, parentIter gtki.TreeIter) {
	iter := r.model.Append(parentIter)

	_ = r.model.SetValue(iter, roomViewRosterStatusIconIndex, getOccupantIconForStatus(o.Status))
	_ = r.model.SetValue(iter, roomViewRosterNicknameIndex, o.Nickname)
	_ = r.model.SetValue(iter, roomViewRosterAffiliationIndex, affiliationDisplayName(o.Affiliation))
	_ = r.model.SetValue(iter, roomViewRosterInfoIndex, occupantDisplayTooltip(o))
}

// parentWindow MUST be called from the UI threads
func (r *roomViewRoster) parentWindow() gtki.Window {
	return r.roomView.mainWindow()
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
