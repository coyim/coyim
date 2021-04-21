package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type occupantActionViewData struct {
	occupant   *muc.Occupant
	rosterView *roomViewRoster
}

type occupantActionView struct {
	occupant   *muc.Occupant
	rosterView *roomViewRoster

	dialog      gtki.Dialog   `gtk-widget:"occupant-action-dialog"`
	header      gtki.Label    `gtk-widget:"occupant-action-header"`
	message     gtki.Label    `gtk-widget:"occupant-action-message"`
	reason      gtki.TextView `gtk-widget:"occupant-action-reason-entry"`
	reasonLabel gtki.Label    `gtk-widget:"occupant-action-reason-label"`

	confirmationAction func() // confirmationAction will not be called from the UI thread
}

func newOccupantActionView(d *occupantActionViewData) *occupantActionView {
	oa := &occupantActionView{
		occupant:   d.occupant,
		rosterView: d.rosterView,
	}

	oa.initBuilder()

	return oa
}

func (oa *occupantActionView) initBuilder() {
	b := newBuilder("MUCRoomOccupantActionDialog")
	panicOnDevError(b.bindObjects(oa))

	b.ConnectSignals(map[string]interface{}{
		"on_ok":     oa.onConfirmClicked,
		"on_cancel": oa.onCancelClicked,
	})
}

func (r *roomViewRoster) newKickOccupantView(o *muc.Occupant) *occupantActionView {
	k := newOccupantActionView(&occupantActionViewData{
		occupant:   o,
		rosterView: r,
	})

	k.initKickOccupantDefaults()

	return k
}

// initKickOccupantDefaults MUST be called from the UI thread
func (oa *occupantActionView) initKickOccupantDefaults() {
	oa.dialog.SetTitle(i18n.Localf("Expel %s from the room", oa.occupant.Nickname))
	oa.header.SetText(i18n.Localf("You are about to temporarily remove %s from the room.", oa.occupant.Nickname))
	oa.message.SetText(i18n.Localf("They will be able to join the room again. Are you sure you want to continue?"))
	oa.reasonLabel.SetText(i18n.Localf("Here you can provide an optional reason for removing the person. Everyone in the room will see this reason."))

	oa.dialog.SetTransientFor(oa.rosterView.roomView.window)
	mucStyles.setRoomDialogErrorComponentHeaderStyle(oa.header)

	oa.confirmationAction = func() {
		oa.rosterView.updateOccupantRole(oa.occupant, &data.NoneRole{}, getTextViewText(oa.reason))
	}
}

func (r *roomViewRoster) newBanOccupantView(o *muc.Occupant) *occupantActionView {
	k := newOccupantActionView(&occupantActionViewData{
		occupant:   o,
		rosterView: r,
	})

	k.initBanOccupantDefaults()

	return k
}

// initBanOccupantDefaults MUST be called from the UI thread
func (oa *occupantActionView) initBanOccupantDefaults() {
	oa.dialog.SetTitle(i18n.Localf("Ban %s from the room", oa.occupant.Nickname))
	oa.header.SetText(i18n.Localf("You are about to ban %s from the room", oa.occupant.Nickname))
	oa.message.SetText(i18n.Local("They won't be able to join the room again. Are you sure you want to continue?"))
	oa.reasonLabel.SetText(i18n.Local("Here you can provide an optional reason for banning the person. Everyone in the room will see this reason."))

	oa.dialog.SetTransientFor(oa.rosterView.roomView.window)
	mucStyles.setRoomDialogErrorComponentHeaderStyle(oa.header)

	oa.confirmationAction = func() {
		oa.rosterView.updateOccupantAffiliation(oa.occupant, &data.OutcastAffiliation{}, getTextViewText(oa.reason))
	}
}

// onConfirmClicked MUST be called from the UI thread
func (oa *occupantActionView) onConfirmClicked() {
	go oa.confirmationAction()
	oa.close()
}

// onCancelClicked MUST be called from the UI thread
func (oa *occupantActionView) onCancelClicked() {
	oa.close()
}

// show MUST be called from the UI thread
func (oa *occupantActionView) show() {
	oa.dialog.Show()
}

// close MUST be called from the UI thread
func (oa *occupantActionView) close() {
	oa.dialog.Destroy()
}
