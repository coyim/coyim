package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type kickOccupantView struct {
	occupant *muc.Occupant

	dialog gtki.Dialog   `gtk-widget:"kick-room-dialog"`
	title  gtki.Label    `gtk-widget:"title-kick-occupant"`
	reason gtki.TextView `gtk-widget:"kick-occupant-reason-entry"`
}

func (v *roomView) newKickOccupantView(o *muc.Occupant) *kickOccupantView {
	d := &kickOccupantView{
		occupant: o,
	}

	d.initBuilder()
	d.initDefaults(v)

	return d
}

func (k *kickOccupantView) initBuilder() {
	b := newBuilder("MUCRoomKickOccupantDialog")
	panicOnDevError(b.bindObjects(k))

	b.ConnectSignals(map[string]interface{}{
		"on_ok":     k.onKickClicked,
		"on_cancel": k.onCancelClicked,
	})
}

func (k *kickOccupantView) initDefaults(v *roomView) {
	k.dialog.SetTransientFor(v.window)
	k.title.SetText(i18n.Localf("You are kicking %s", k.occupant.Nickname))
}

// close MUST be called from the UI thread
func (k *kickOccupantView) onKickClicked() {
	// TODO: Implement
	k.close()
}

// close MUST be called from the UI thread
func (k *kickOccupantView) onCancelClicked() {
	k.close()
}

func (k *kickOccupantView) show() {
	k.dialog.Show()
}

func (k *kickOccupantView) close() {
	k.dialog.Destroy()
}
