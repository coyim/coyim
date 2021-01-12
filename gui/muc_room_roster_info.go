package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomViewRosterInfo struct {
	u *gtkUI

	account  *account
	roomID   jid.Bare
	occupant *muc.Occupant

	view                    gtki.Box   `gtk-widget:"roster-info-box"`
	avatar                  gtki.Image `gtk-widget:"occupant-avatar"`
	nicknameLabel           gtki.Label `gtk-widget:"occupant-nickname"`
	realJIDLabel            gtki.Label `gtk-widget:"user-jid"`
	status                  gtki.Label `gtk-widget:"status"`
	statusMessage           gtki.Label `gtk-widget:"status-message"`
	currentAffiliationLabel gtki.Label `gtk-widget:"current-affiliation"`

	onReset     *callbacksSet
	onRefresh   *callbacksSet
	onHidePanel func()

	log coylog.Logger
}

func (r *roomViewRoster) newRoomViewRosterInfo(onHidePanel func()) *roomViewRosterInfo {
	ri := &roomViewRosterInfo{
		u:           r.u,
		account:     r.accout,
		roomID:      r.roomID,
		onReset:     newCallbacksSet(),
		onRefresh:   newCallbacksSet(),
		onHidePanel: onHidePanel,
		log:         r.log,
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
	)

	r.onReset.add(
		r.removeOccupantInfo,
		r.removeOccupantAffiliationInfo,
	)
}

func (r *roomViewRosterInfo) occupantAffiliationChanged() {
	r.log.WithFields(log.Fields{
		"where":       "occupantAffiliationUpdate",
		"occupant":    r.occupant.RealJid,
		"affiliation": r.occupant.Affiliation.Name(),
	}).Info("The occupant affiliation has been updated")

	doInUIThread(r.refresh)
}

// showOccupantInfo MUST be called from the UI thread
func (r *roomViewRosterInfo) showOccupantInfo(occupant *muc.Occupant) {
	r.occupant = occupant
	r.refresh()
	r.show()
}

// refresh MUST be called from the UI thread
func (r *roomViewRosterInfo) refresh() {
	r.reset()
	if r.account != nil {
		r.onRefresh.invokeAll()
	}
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

// show MUST be called from the UI thread
func (r *roomViewRosterInfo) show() {
	r.view.Show()
}

// show MUST be called from the UI thread
func (r *roomViewRosterInfo) hide() {
	r.view.Hide()

	if r.onHidePanel != nil {
		r.onHidePanel()
	}

	r.reset()
}

// widget MUST be called from the UI thread
func (r *roomViewRosterInfo) widget() gtki.Widget {
	return r.view
}
