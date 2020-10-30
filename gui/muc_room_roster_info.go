package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewRosterInfo struct {
	rosterInfoBox  gtki.Box      `gtk-widget:"roster-info-box"`
	rosterInfo     gtki.Revealer `gtk-widget:"roster-info-revelear"`
	occupantAvatar gtki.Image    `gtk-widget:"occupant-avatar"`
	nickname       gtki.Label    `gtk-widget:"occupant-nickname"`
	userJID        gtki.Label    `gtk-widget:"user-jid"`
	status         gtki.Label    `gtk-widget:"status"`
	statusMessage  gtki.Label    `gtk-widget:"status-message"`
	role           gtki.Label    `gtk-widget:"role"`
	affiliation    gtki.Label    `gtk-widget:"affiliation"`
	voice          gtki.Label    `gtk-widget:"voice"`
	onHidePanel    func()
}

func (v *roomView) newRoomViewRosterInfo() *roomViewRosterInfo {
	r := &roomViewRosterInfo{}
	r.initBuilder()
	r.initDefaults()

	return r
}

func (r *roomViewRosterInfo) initBuilder() {
	builder := newBuilder("MUCRoomRosterInfo")
	builder.ConnectSignals(map[string]interface{}{
		"on_hide": r.onHideOccupantInfoPanel,
	})

	panicOnDevError(builder.bindObjects(r))

}

func (r *roomViewRosterInfo) initDefaults() {
	r.rosterInfo.Hide()
}

func (r *roomViewRosterInfo) onHideOccupantInfoPanel() {
	r.rosterInfo.Hide()
	r.onHidePanel()
}

func (r *roomViewRosterInfo) displayOccupantInfoPanel(occupant *muc.Occupant, showRoster func()) {
	r.onHidePanel = showRoster

	r.populateOccupantInfoPanel(occupant)

	r.rosterInfo.SetRevealChild(true)
	r.rosterInfo.Show()
}

func (r *roomViewRosterInfo) populateOccupantInfoPanel(occupant *muc.Occupant) {
	r.occupantAvatar.SetFromPixbuf(getMUCIconPixbuf(getOccupantIconNameForStatus(occupant.Status.Status)))

	r.nickname.SetText(i18n.Localf("About: %s", occupant.Nickname))

	r.userJID.SetText("Undefined")
	if occupant.RealJid != nil {
		rj := occupant.RealJid.String()
		r.userJID.SetText(rj)
		r.userJID.SetTooltipText(rj)
	}

	r.status.SetText(showForDisplay(occupant.Status.Status, false))

	r.statusMessage.SetText("-")
	if occupant.Status.StatusMsg != "" {
		sm := occupant.Status.StatusMsg
		r.statusMessage.SetText(sm)
		r.statusMessage.SetTooltipText(sm)
	}

	r.role.SetText(roleDisplayName(occupant.Role))
	r.affiliation.SetText(affiliationDisplayName(occupant.Affiliation))

	r.voice.SetText("No")
	if occupant.HasVoice() {
		r.voice.SetText("Yes")
	}
}
