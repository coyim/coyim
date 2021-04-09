package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type occupantActionView struct {
	occupant       *muc.Occupant
	roomViewRoster *roomViewRoster

	dialog      gtki.Dialog   `gtk-widget:"occupant-action-dialog"`
	header      gtki.Label    `gtk-widget:"occupant-action-header"`
	message     gtki.Label    `gtk-widget:"occupant-action-message"`
	reason      gtki.TextView `gtk-widget:"occupant-action-reason-entry"`
	reasonLabel gtki.Label    `gtk-widget:"occupant-action-reason-label"`
}

func (r *roomViewRoster) newOccupantActionView(o *muc.Occupant) *occupantActionView {
	oa := &occupantActionView{
		occupant:       o,
		roomViewRoster: r,
	}

	oa.initBuilder()
	return oa
}

func (r *roomViewRoster) newKickOccupantView(o *muc.Occupant) *occupantActionView {
	k := r.newOccupantActionView(o)
	k.initKickOccupantDefaults()

	return k
}

func (r *roomViewRoster) newChangeOccupantVoiceView(o *muc.Occupant) *occupantActionView {
	k := r.newOccupantActionView(o)
	k.initChangeOccupantVoiceDefaults()

	return k
}

func (oa *occupantActionView) initBuilder() {
	b := newBuilder("MUCRoomOccupantActionDialog")
	panicOnDevError(b.bindObjects(oa))

	b.ConnectSignals(map[string]interface{}{
		"on_ok":     oa.onKickClicked,
		"on_cancel": oa.onCancelClicked,
	})
}

// initDefaults MUST be called from the UI thread
func (oa *occupantActionView) initKickOccupantDefaults() {
	oa.dialog.SetTitle(i18n.Localf("Expel %s from the room", oa.occupant.Nickname))
	oa.dialog.SetTransientFor(oa.roomViewRoster.roomView.window)

	oa.header.SetText(i18n.Localf("You are about to temporarily remove %s from the room.", oa.occupant.Nickname))
	oa.message.SetText(i18n.Localf("They will be able to join the room again. Are you sure you want to continue?"))
	oa.reasonLabel.SetText(i18n.Localf("Here you can provide an optional reason for removing the person.\nEveryone in the room will see this reason."))
	mucStyles.setRoomDialogErrorComponentHeaderStyle(oa.header)
}

// initDefaults MUST be called from the UI thread
func (oa *occupantActionView) initChangeOccupantVoiceDefaults() {
	title := i18n.Localf("Grant voice to %s", oa.occupant.Nickname)
	header := i18n.Localf("You are about to change the role of %s from visitor to participant", oa.occupant.Nickname)
	message := i18n.Local("It allows to send messages in the room. Are you sure you want to continue?")
	if oa.occupant.HasVoice() {
		title = i18n.Localf("Revoke voice to %s", oa.occupant.Nickname)
		header = i18n.Localf("You are about to change the role of %s to visitor", oa.occupant.Nickname)
		message = i18n.Local("A visitor is not able to send messages in the room. Are you sure you want to continue?")
	}
	oa.dialog.SetTitle(title)
	oa.header.SetText(header)
	oa.message.SetText(message)
	oa.reasonLabel.SetText(i18n.Local("Here you can provide an optional reason. Everyone in the room will see it."))

	oa.dialog.SetTransientFor(oa.roomViewRoster.roomView.window)
	mucStyles.setRoomDialogErrorComponentHeaderStyle(oa.header)
}

// onKickClicked MUST be called from the UI thread
func (oa *occupantActionView) onKickClicked() {
	go oa.roomViewRoster.kickOccupant(oa.occupant, getTextViewText(oa.reason))
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
