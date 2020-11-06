package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomDestroyView struct {
	room *roomView

	transient           gtki.Window
	dialog              gtki.Dialog `gtk-widget:"destroy-room-dialog"`
	reasonEntry         gtki.Entry  `gtk-widget:"destroy-room-reason-entry"`
	alternateVenueEntry gtki.Entry  `gtk-widget:"destroy-room-alternate-venue-entry"`
	destroyRoomButton   gtki.Button `gtk-widget:"destroy-room-button"`
	spinnerBox          gtki.Box    `gtk-widget:"destroy-room-spinner-box"`
	notificationBox     gtki.Box    `gtk-widget:"notification-area"`

	spinner      *spinner
	notification *notifications

	cancelChannel chan bool
}

func (v *roomView) newRoomDestroyView(t gtki.Window) *roomDestroyView {
	d := &roomDestroyView{
		room:      v,
		transient: t,
	}

	d.initBuilder()
	d.initDefaults(v.u)

	return d
}

func (d *roomDestroyView) initBuilder() {
	builder := newBuilder("MUCRoomDestroyDialog")
	panicOnDevError(builder.bindObjects(d))

	builder.ConnectSignals(map[string]interface{}{
		"on_destroy_clicked":  d.onDestroyRoom,
		"on_cancel_clicked":   d.onCancel,
		"on_dialog_destroyed": d.onDialogDestroy,
	})
}

func (d *roomDestroyView) initDefaults(u *gtkUI) {
	d.dialog.SetTransientFor(d.transient)

	d.spinner = newSpinner()
	d.spinnerBox.Add(d.spinner.getWidget())

	d.notification = u.newNotifications(d.notificationBox)
}

func (d *roomDestroyView) onDestroyRoom() {
	d.spinner.show()
	d.disableFields()

	reason := d.getReason()
	alternateID := d.getAlternateID()

	d.room.tryDestroyRoom(alternateID, reason, d.onDestroySuccess, d.onDestroyFails)
}

func (d *roomDestroyView) onDestroySuccess() {
	d.close()
}

func (d *roomDestroyView) onDestroyFails() {
	d.enableFields()
	// TODO: check which type of error happened to show a friendly message to the user
	d.notification.error(i18n.Local("The room can't be destroyed, please try again."))
}

func (d *roomDestroyView) onCancel() {
	d.spinner.hide()
	d.cancelActiveRequest()
	d.close()
}

func (d *roomDestroyView) onDialogDestroy() {
	d.cancelActiveRequest()
}

func (d *roomDestroyView) cancelActiveRequest() {
	if d.cancelChannel != nil {
		d.cancelChannel <- true
	}
}

func (d *roomDestroyView) getReason() string {
	t, _ := d.reasonEntry.GetText()
	return t
}

func (d *roomDestroyView) getAlternateID() jid.Bare {
	t, _ := d.alternateVenueEntry.GetText()
	if t != "" {
		return jid.ParseBare(t)
	}
	return nil
}

func (d *roomDestroyView) disableFields() {
	d.reasonEntry.SetSensitive(false)
	d.alternateVenueEntry.SetSensitive(false)
	d.destroyRoomButton.SetSensitive(false)
}

func (d *roomDestroyView) enableFields() {
	d.reasonEntry.SetSensitive(true)
	d.alternateVenueEntry.SetSensitive(true)
	d.destroyRoomButton.SetSensitive(true)
}

func (d *roomDestroyView) show() {
	d.dialog.Show()
}

func (d *roomDestroyView) close() {
	d.dialog.Destroy()
}
