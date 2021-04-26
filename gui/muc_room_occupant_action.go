package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type occupantActionViewData struct {
	parentWindow gtki.Window

	dialogTitle string
	headerText  string
	messageText string
	reasonText  string

	confirmationAction func(reason string) // confirmationAction will not be called from the UI thread
}

type occupantActionView struct {
	dialog      gtki.Dialog   `gtk-widget:"occupant-action-dialog"`
	header      gtki.Label    `gtk-widget:"occupant-action-header"`
	message     gtki.Label    `gtk-widget:"occupant-action-message"`
	reason      gtki.TextView `gtk-widget:"occupant-action-reason-entry"`
	reasonLabel gtki.Label    `gtk-widget:"occupant-action-reason-label"`

	confirmationAction func(reason string) // confirmationAction will not be called from the UI thread
}

func newOccupantActionView(d *occupantActionViewData) *occupantActionView {
	oa := &occupantActionView{
		confirmationAction: d.confirmationAction,
	}

	oa.initBuilder()
	oa.initDefaults(d)
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

func (oa *occupantActionView) initDefaults(d *occupantActionViewData) {
	oa.dialog.SetTransientFor(d.parentWindow)
	mucStyles.setRoomDialogErrorComponentHeaderStyle(oa.header)
}

func (oa *occupantActionView) initDialogTitleAndTexts(d *occupantActionViewData) {
	oa.dialog.SetTitle(d.dialogTitle)
	oa.header.SetText(d.headerText)
	oa.message.SetText(d.messageText)
	oa.reasonLabel.SetText(d.reasonText)
}

// onConfirmClicked MUST be called from the UI thread
func (oa *occupantActionView) onConfirmClicked() {
	reason := getTextViewText(oa.reason)
	go oa.confirmationAction(reason)
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

func (r *roomViewRoster) newKickOccupantView(o *muc.Occupant) *occupantActionView {
	k := newOccupantActionView(&occupantActionViewData{
		parentWindow: r.parentWindow(),

		dialogTitle: i18n.Localf("Expel %s from the room", o.Nickname),
		headerText:  i18n.Localf("You are about to temporarily remove %s from the room.", o.Nickname),
		messageText: i18n.Local("They will be able to join the room again. Are you sure you want to continue?"),
		reasonText:  i18n.Local("Here you can provide an optional reason for removing the person. Everyone in the room will see this reason."),

		confirmationAction: func(reason string) {
			r.updateOccupantRole(o, &data.NoneRole{}, reason)
			r.hideRosterInfoPanel()
		},
	})

	return k
}

func (r *roomViewRoster) newBanOccupantView(o *muc.Occupant) *occupantActionView {
	k := newOccupantActionView(&occupantActionViewData{
		parentWindow: r.parentWindow(),

		dialogTitle: i18n.Localf("Ban %s from the room", o.Nickname),
		headerText:  i18n.Localf("You are about to ban %s from the room", o.Nickname),
		messageText: i18n.Local("They won't be able to join the room again. Are you sure you want to continue?"),
		reasonText:  i18n.Local("Here you can provide an optional reason for banning the person. Everyone in the room will see this reason."),

		confirmationAction: func(reason string) {
			r.updateOccupantAffiliation(o, &data.OutcastAffiliation{}, reason)
			r.hideRosterInfoPanel()
		},
	})

	return k
}
