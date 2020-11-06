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
	onHidePanel    func()
}

func (r *roomViewRoster) newRoomViewRosterInfo() *roomViewRosterInfo {
	ri := &roomViewRosterInfo{}
	ri.initBuilder()
	ri.initDefaults()

	return ri
}

func (r *roomViewRosterInfo) initBuilder() {
	builder := newBuilder("MUCRoomRosterInfo")
	builder.ConnectSignals(map[string]interface{}{
		"on_hide": r.onHideOccupantInfoPanel,
	})

	panicOnDevError(builder.bindObjects(r))

}

func (r *roomViewRosterInfo) initDefaults() {
	updateWithStyle(r.nickname, providerWithStyle("label", style{
		"font-size":   "14px",
		"font-weight": "bold",
	}))

	updateWithStyle(r.userJID, providerWithStyle("label", style{
		"font-size": "12px",
	}))

	updateWithStyle(r.status, providerWithStyle("label", style{
		"font-size":   "12px",
		"font-style":  "italic",
		"font-weight": "bold",
		"color":       "#666666",
	}))
}

func (r *roomViewRosterInfo) onHideOccupantInfoPanel() {
	r.rosterInfo.Hide()
	r.onHidePanel()
}

func (r *roomViewRosterInfo) displayOccupantInfoPanel(occupant *muc.Occupant, showRoster func()) {
	r.onHidePanel = showRoster

	r.populateOccupantInfoPanel(occupant)
	r.rosterInfo.Show()
}

func (r *roomViewRosterInfo) populateOccupantInfoPanel(occupant *muc.Occupant) {
	r.occupantAvatar.SetFromPixbuf(getMUCIconPixbuf(getOccupantIconNameForStatus(occupant.Status.Status)))

	r.nickname.SetText(i18n.Local(occupant.Nickname))

	r.userJID.SetVisible(false)
	if occupant.RealJid != nil {
		rj := occupant.RealJid.String()
		r.userJID.SetText(rj)
		r.userJID.SetTooltipText(rj)
		r.userJID.SetVisible(true)
	}

	r.status.SetText(showForDisplay(occupant.Status.Status, false))

	if occupant.Status.StatusMsg != "" {
		sm := occupant.Status.StatusMsg
		r.statusMessage.SetText(sm)
		r.statusMessage.SetTooltipText(sm)
	}
}
