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

	dialogTitle string
	headerText  string
	messageText string
	reasonText  string
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
	oa.initDefaults()
	oa.initDialogTitleAndTexts(d)

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

func (oa *occupantActionView) initDefaults() {
	oa.dialog.SetTransientFor(oa.rosterView.parentWindow())
	mucStyles.setRoomDialogErrorComponentHeaderStyle(oa.header)
}

func (oa *occupantActionView) initDialogTitleAndTexts(d *occupantActionViewData) {
	oa.dialog.SetTitle(d.dialogTitle)
	oa.header.SetText(d.headerText)
	oa.message.SetText(d.messageText)
	oa.reasonLabel.SetText(d.reasonText)
}

func (r *roomViewRoster) newKickOccupantView(o *muc.Occupant) *occupantActionView {
	k := newOccupantActionView(&occupantActionViewData{
		occupant:   o,
		rosterView: r,

		dialogTitle: i18n.Localf("Expel %s from the room", o.Nickname),
		headerText:  i18n.Localf("You are about to temporarily remove %s from the room.", o.Nickname),
		messageText: i18n.Local("They will be able to join the room again. Are you sure you want to continue?"),
		reasonText:  i18n.Local("Here you can provide an optional reason for removing the person. Everyone in the room will see this reason."),
	})

	k.initKickOccupantDefaults()

	return k
}

// initKickOccupantDefaults MUST be called from the UI thread
func (oa *occupantActionView) initKickOccupantDefaults() {
	oa.confirmationAction = func() {
		oa.rosterView.updateOccupantRole(oa.occupant, &data.NoneRole{}, getTextViewText(oa.reason))
	}
}

func (r *roomViewRoster) newBanOccupantView(o *muc.Occupant) *occupantActionView {
	k := newOccupantActionView(&occupantActionViewData{
		occupant:   o,
		rosterView: r,

		dialogTitle: i18n.Localf("Ban %s from the room", o.Nickname),
		headerText:  i18n.Localf("You are about to ban %s from the room", o.Nickname),
		messageText: i18n.Local("They won't be able to join the room again. Are you sure you want to continue?"),
		reasonText:  i18n.Local("Here you can provide an optional reason for banning the person. Everyone in the room will see this reason."),
	})

	k.initBanOccupantDefaults()

	return k
}

// initBanOccupantDefaults MUST be called from the UI thread
func (oa *occupantActionView) initBanOccupantDefaults() {
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
