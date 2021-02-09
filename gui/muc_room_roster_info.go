package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewRosterInfo struct {
	u *gtkUI

	account      *account
	roomID       jid.Bare
	occupant     *muc.Occupant
	selfOccupant *muc.Occupant
	rosterView   *roomViewRoster

	view                    gtki.Box    `gtk-widget:"roster-info-box"`
	avatar                  gtki.Image  `gtk-widget:"occupant-avatar"`
	nicknameLabel           gtki.Label  `gtk-widget:"occupant-nickname"`
	realJIDLabel            gtki.Label  `gtk-widget:"user-jid"`
	status                  gtki.Label  `gtk-widget:"status"`
	statusMessage           gtki.Label  `gtk-widget:"status-message"`
	currentAffiliationLabel gtki.Label  `gtk-widget:"current-affiliation"`
	currentRoleLabel        gtki.Label  `gtk-widget:"current-role"`
	changeRoleButton        gtki.Button `gtk-widget:"change-role"`
	changeAffiliationButton gtki.Button `gtk-widget:"change-affiliation"`

	onReset   *callbacksSet
	onRefresh *callbacksSet

	log coylog.Logger
}

func (r *roomViewRoster) newRoomViewRosterInfo() *roomViewRosterInfo {
	ri := &roomViewRosterInfo{
		u:          r.u,
		account:    r.account,
		roomID:     r.roomID,
		rosterView: r,
		onReset:    newCallbacksSet(),
		onRefresh:  newCallbacksSet(),
		log:        r.log,
	}

	ri.initBuilder()
	ri.initCSSStyles()
	ri.initDefaults()

	return ri
}

func (r *roomViewRosterInfo) initBuilder() {
	builder := newBuilder("MUCRoomRosterInfo")
	panicOnDevError(builder.bindObjects(r))

	builder.ConnectSignals(map[string]interface{}{
		"on_hide":               r.hide,
		"on_change_affiliation": r.onChangeAffiliation,
		"on_change_role":        r.onChangeRole,
	})
}

func (r *roomViewRosterInfo) initCSSStyles() {
	mucStyles.setRoomRosterInfoNicknameLabelStyle(r.nicknameLabel)
	mucStyles.setRoomRosterInfoUserJIDLabelStyle(r.realJIDLabel)
	mucStyles.setRoomRosterInfoStatusLabelStyle(r.status)
}

func (r *roomViewRosterInfo) initDefaults() {
	r.onRefresh.add(
		r.refreshOccupantInfo,
		r.refreshOccupantAffiliation,
		r.refreshOccupantRole,
		r.validateOccupantPrivileges,
	)

	r.onReset.add(
		r.removeOccupantInfo,
		r.removeOccupantAffiliationInfo,
		r.removeOccupantRoleInfo,
		r.validateOccupantPrivileges,
	)
}

func (r *roomViewRosterInfo) updateOccupantAffiliation(occupant *muc.Occupant, previousAffiliation data.Affiliation, reason string) {
	r.rosterView.updateOccupantAffiliation(occupant, previousAffiliation, reason)
	doInUIThread(r.refresh)
}

func (r *roomViewRosterInfo) updateOccupantRole(occupant *muc.Occupant, role data.Role, reason string) {
	r.rosterView.updateOccupantRole(occupant, role, reason)
	doInUIThread(r.refresh)
}

func (r *roomViewRosterInfo) updateSelfOccupant(o *muc.Occupant) {
	r.selfOccupant = o
}

// showOccupantInfo MUST be called from the UI thread
func (r *roomViewRosterInfo) showOccupantInfo(occupant *muc.Occupant) {
	r.occupant = occupant
	r.refresh()
	r.show()
}

func (r *roomViewRosterInfo) validateOccupantPrivileges() {
	// Privileges associated with affiliations
	// See: https://xmpp.org/extensions/xep-0045.html#affil-priv
	switch r.selfOccupant.Affiliation.(type) {
	case *data.OwnerAffiliation:
		r.changeAffiliationButton.SetVisible(true)
	case *data.AdminAffiliation:
		r.changeAffiliationButton.SetVisible(r.occupant.Affiliation.IsLowerThan(r.selfOccupant.Affiliation))
	default:
		r.changeAffiliationButton.SetVisible(false)
	}

	canChangeRole := (r.selfOccupant.Affiliation.IsAdmin() || r.selfOccupant.Affiliation.IsOwner()) &&
		!r.occupant.Affiliation.IsAdmin() && !r.occupant.Affiliation.IsOwner()
	r.changeRoleButton.SetVisible(canChangeRole)
}

// refresh MUST be called from the UI thread
func (r *roomViewRosterInfo) refresh() {
	r.reset()
	r.onRefresh.invokeAll()
}

// reset MUST be called from the UI thread
func (r *roomViewRosterInfo) reset() {
	r.onReset.invokeAll()
}

// refresh MUST be called from the UI thread
func (r *roomViewRosterInfo) refreshOccupantInfo() {
	occupant := r.occupant
	status := r.occupant.Status

	r.avatar.SetFromPixbuf(getMUCIconPixbuf(getOccupantIconNameForStatus(status.Status)))
	setLabelText(r.nicknameLabel, occupant.Nickname)

	if occupant.RealJid != nil {
		r.realJIDLabel.SetText(occupant.RealJid.String())
		r.realJIDLabel.SetTooltipText(occupant.RealJid.String())
		r.realJIDLabel.SetVisible(true)
	}

	r.status.SetText(showForDisplay(status.Status, false))
	if status.StatusMsg != "" {
		r.statusMessage.SetText(status.StatusMsg)
		r.statusMessage.SetTooltipText(status.StatusMsg)
		r.statusMessage.SetVisible(true)
	}
}

// removeOccupantInfo MUST be called from the UI thread
func (r *roomViewRosterInfo) removeOccupantInfo() {
	r.avatar.Clear()

	r.nicknameLabel.SetText("")

	r.realJIDLabel.SetText("")
	r.realJIDLabel.SetVisible(false)

	r.statusMessage.SetText("")
	r.statusMessage.SetVisible(false)
}

// refreshOccupantAffiliation MUST be called from the UI thread
func (r *roomViewRosterInfo) refreshOccupantAffiliation() {
	r.currentAffiliationLabel.SetText(occupantAffiliationName(r.occupant.Affiliation))
}

// removeOccupantAffiliationInfo MUST be called from the UI thread
func (r *roomViewRosterInfo) removeOccupantAffiliationInfo() {
	r.currentAffiliationLabel.SetText("")
}

// refreshOccupantAffiliation MUST be called from the UI thread
func (r *roomViewRosterInfo) refreshOccupantRole() {
	r.currentRoleLabel.SetText(occupantRoleName(r.occupant.Role))
}

// removeOccupantAffiliationInfo MUST be called from the UI thread
func (r *roomViewRosterInfo) removeOccupantRoleInfo() {
	r.currentRoleLabel.SetText("")
}

// show MUST be called from the UI thread
func (r *roomViewRosterInfo) show() {
	r.view.Show()
}

// show MUST be called from the UI thread
func (r *roomViewRosterInfo) hide() {
	r.view.Hide()
	r.rosterView.hideRosterInfoPanel()
	r.reset()
}

// parentWindow MUST be called from the UI threads
func (r *roomViewRosterInfo) parentWindow() gtki.Window {
	return r.rosterView.parentWindow()
}

func (r *roomViewRosterInfo) contentBox() gtki.Box {
	return r.view
}
