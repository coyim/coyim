package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type kickOccupantView struct {
	occupant       *muc.Occupant
	roomViewRoster *roomViewRoster

	dialog gtki.Dialog   `gtk-widget:"kick-room-dialog"`
	title  gtki.Label    `gtk-widget:"title-kick-occupant"`
	reason gtki.TextView `gtk-widget:"kick-occupant-reason-entry"`
}

func (r *roomViewRoster) newKickOccupantView(o *muc.Occupant) *kickOccupantView {
	k := &kickOccupantView{
		occupant:       o,
		roomViewRoster: r,
	}

	k.initBuilder()
	k.initDefaults()

	return k
}

func (k *kickOccupantView) initBuilder() {
	b := newBuilder("MUCRoomKickOccupantDialog")
	panicOnDevError(b.bindObjects(k))

	b.ConnectSignals(map[string]interface{}{
		"on_ok":     k.onKickClicked,
		"on_cancel": k.onCancelClicked,
	})
}

// initDefaults MUST be called from the UI thread
func (k *kickOccupantView) initDefaults() {
	k.dialog.SetTitle(i18n.Localf("Expel %s from the room", k.occupant.Nickname))
	k.dialog.SetTransientFor(k.roomViewRoster.roomView.window)

	k.title.SetText(i18n.Localf("You are about to temporarily remove %s from the room.", k.occupant.Nickname))
	mucStyles.setRoomDialogErrorComponentHeaderStyle(k.title)
}

// onKickClicked MUST be called from the UI thread
func (k *kickOccupantView) onKickClicked() {
	go k.roomViewRoster.updateOccupantRole(k.occupant, &data.NoneRole{}, getTextViewText(k.reason))
	k.close()
}

// onCancelClicked MUST be called from the UI thread
func (k *kickOccupantView) onCancelClicked() {
	k.close()
}

// show MUST be called from the UI thread
func (k *kickOccupantView) show() {
	k.dialog.Show()
}

// close MUST be called from the UI thread
func (k *kickOccupantView) close() {
	k.dialog.Destroy()
}
