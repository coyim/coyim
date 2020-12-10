package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomDestroyErrorView struct {
	roomID       jid.Bare
	destroyError error
	onRetry      func()

	dialog       gtki.Dialog `gtk-widget:"destroy-room-error-dialog"`
	errorMessage gtki.Label  `gtk-widget:"destroy-room-error-message"`
}

func newDestroyError(roomID jid.Bare, err error, onRetry func()) *roomDestroyErrorView {
	rd := &roomDestroyErrorView{
		roomID:       roomID,
		destroyError: err,
		onRetry:      onRetry,
	}

	rd.initBuilder()
	rd.initDefaults()

	return rd
}

func (rd *roomDestroyErrorView) initBuilder() {
	b := newBuilder("MUCRoomDestroyErrorDialog")
	panicOnDevError(b.bindObjects(rd))

	b.ConnectSignals(map[string]interface{}{
		"on_retry":  rd.retry,
		"on_cancel": rd.close,
	})
}

func (rd *roomDestroyErrorView) initDefaults() {
	rd.dialog.SetTitle(i18n.Localf("Room [%s] destroy error", rd.roomID))
	rd.errorMessage.SetLabel(rd.friendlyMessageForDestroyError(rd.destroyError))
}

func (rd *roomDestroyErrorView) retry() {
	rd.onRetry()
	rd.close()
}

func (rd *roomDestroyErrorView) close() {
	rd.dialog.Destroy()
}

func (rd *roomDestroyErrorView) show() {
	rd.dialog.Show()
}

func (rd *roomDestroyErrorView) friendlyMessageForDestroyError(err error) string {
	switch err {
	case session.ErrDestroyRoomInvalidIQResponse:
		return i18n.Local("We were able to connect to the room service, " +
			"but we received an invalid response from it. Please try again.")
	case session.ErrDestroyRoomForbidden:
		return i18n.Local("You don't have the permission to destroy this room. " +
			"Please get in contact with one of the room owners.")
	case session.ErrDestroyRoomUnknown:
		return i18n.Local("The room's service responded with an unknow error, " +
			"so, the room can't be destroyed. Please try again.")
	case session.ErrDestroyRoomNoResult:
		return i18n.Local("We were able to send the request to destroy the room, " +
			"but the service responded with an unknow result. Please contact the " +
			"room's administrator.")
	default:
		return i18n.Local("An error occurred when destroying the room please try again.")
	}
}
