package gui

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/coylog"
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

var (
	errCreateRoomCheckIfExistsFails = errors.New("room exists failed")
	errCreateRoomAlreadyExists      = errors.New("room already exists")
	errCreateRoomFailed             = errors.New("couldn't create the room")
)

type mucCreateRoomView struct {
	u *gtkUI

	autoJoin      bool
	configureRoom bool
	cancel        chan bool

	dialog    gtki.Dialog `gtk-widget:"create-room-dialog"`
	container gtki.Box    `gtk-widget:"create-room-content"`

	form    *mucCreateRoomViewForm
	success *mucCreateRoomViewSuccess

	onCreateOptionChange *callbacksSet
	onDestroy            *callbacksSet

	sync.Mutex
}

func newCreateMUCRoomView(u *gtkUI, d *mucCreateRoomData) *mucCreateRoomView {
	v := &mucCreateRoomView{
		u:                    u,
		onCreateOptionChange: newCallbacksSet(),
		onDestroy:            newCallbacksSet(),
	}

	v.initBuilder()
	v.initCreateRoomForm(d)
	v.initCreateRoomSuccess()

	return v
}

func (v *mucCreateRoomView) initBuilder() {
	builder := newBuilder("MUCCreateRoomDialog")
	panicOnDevError(builder.bindObjects(v))

	builder.ConnectSignals(map[string]interface{}{
		"on_close_window": v.onCloseWindow,
	})
}

// onCloseWindow MUST be called from the UI thread
func (v *mucCreateRoomView) onCloseWindow() {
	v.onDestroy.invokeAll()
}

// onCancel MUST be called from the UI thread
func (v *mucCreateRoomView) onCancel() {
	if v.cancel != nil {
		v.cancel <- true
		v.cancel = nil
	}

	v.destroy()
}

// destroy MUST be called from the UI thread
func (v *mucCreateRoomView) destroy() {
	v.dialog.Destroy()
}

// createRoom IS SAFE to be called from the UI thread
func (v *mucCreateRoomView) createRoom(ca *account, roomID jid.Bare, onError func(err error)) {
	v.cancel = make(chan bool)

	sc := make(chan bool)
	ec := make(chan error)

	onErrorFinal := onError
	onError = func(err error) {
		if onErrorFinal != nil {
			onErrorFinal(err)
		}
	}

	go func() {
		defer func() {
			v.cancel = nil
		}()

		v.checkIfRoomExists(ca, roomID, sc, ec)

		select {
		case <-sc:
			if v.configureRoom {
				v.createReservedRoom(ca, roomID, onError)
			} else {
				v.createInstantRoom(ca, roomID, onError)
			}
		case err := <-ec:
			onError(err)
		case <-v.cancel:
		}
	}()

}

// joinRoom MUST be called from the UI thread
func (v *mucCreateRoomView) joinRoom(ca *account, roomID jid.Bare) {
	v.dialog.Destroy()
	v.u.joinRoom(ca, roomID, nil)
}

// updateAutoJoinValue IS SAFE to be called from the UI thread
func (v *mucCreateRoomView) updateAutoJoinValue(f bool) {
	v.updateCreateOption("autoJoin", f)
}

// updateConfigureRoomValue IS SAFE to be called from the UI thread
func (v *mucCreateRoomView) updateConfigureRoomValue(f bool) {
	v.updateCreateOption("configRoom", f)
}

// updateCreateOption IS SAFE to be called from the UI thread
func (v *mucCreateRoomView) updateCreateOption(o string, f bool) {
	v.Lock()
	defer v.Unlock()

	previousValue := false

	switch o {
	case "autoJoin":
		previousValue = v.autoJoin
		v.autoJoin = f
	case "configRoom":
		previousValue = v.configureRoom
		v.configureRoom = f
	}

	if previousValue != f {
		v.onCreateOptionChange.invokeAll()
	}
}

// show MUST be called from the UI thread
func (v *mucCreateRoomView) show() {
	v.container.ShowAll()
	v.dialog.Show()
}

// log IS SAFE to be called from the UI thread
func (v *mucCreateRoomView) log(ca *account, roomID jid.Bare) coylog.Logger {
	return ca.log.WithFields(log.Fields{
		"room":  roomID,
		"where": "mucCreateRoomView",
	})
}

type mucCreateRoomData struct {
	ca           *account
	roomName     jid.Local
	where        jid.Domain
	autoJoin     bool
	customConfig bool
}

func (u *gtkUI) mucShowCreateRoomWithData(d *mucCreateRoomData, onViewCreated func(*mucCreateRoomView)) {
	v := newCreateMUCRoomView(u, d)
	u.connectShortcutsChildWindow(v.dialog)

	if onViewCreated != nil {
		onViewCreated(v)
	}

	v.dialog.SetTransientFor(u.window)
	v.show()
}

func (u *gtkUI) mucShowCreateRoomForm(d *mucCreateRoomData) {
	u.mucShowCreateRoomWithData(d, func(v *mucCreateRoomView) {
		v.showCreateForm()
	})
}

func (u *gtkUI) mucShowCreateRoomSuccess(ca *account, roomID jid.Bare, d *mucCreateRoomData) {
	u.mucShowCreateRoomWithData(d, func(v *mucCreateRoomView) {
		v.showSuccessView(ca, roomID)
	})
}

func (u *gtkUI) mucShowCreateRoom() {
	u.mucShowCreateRoomForm(nil)
}
