package gui

import (
	"errors"
	"sync"
	"time"

	"github.com/coyim/coyim/coylog"
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

var (
	errCreateRoomCheckIfExistsFails = errors.New("room exists failed")
	errCreateRoomAlreadyExists      = errors.New("room already exists")
	errCreateRoomFailed             = errors.New("couldn't create the room")
	errCreateRoomTimeout            = errors.New("timeout trying to create the room")
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
	v.initCreateRoomSuccess(d)

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

	onErrorFinal := func(err error) {
		if onError != nil {
			onError(err)
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
				v.instantiatePersistentRoom(ca, roomID, onErrorFinal)
			} else {
				v.createInstantRoom(ca, roomID, onErrorFinal)
			}
		case err := <-ec:
			onError(err)
		case <-time.After(timeoutThreshold):
			onError(errCreateRoomTimeout)
		case <-v.cancel:
		}
	}()

}

// instantiatePersistentRoom IS SAFE to be called from the UI thread
func (v *mucCreateRoomView) instantiatePersistentRoom(ca *account, roomID jid.Bare, onError func(error)) {
	fc, ec := ca.session.CreateReservedRoom(roomID)
	go func() {
		select {
		case err := <-ec:
			onError(err)
		case form := <-fc:
			form.ConfigureRoomAsPersistent()
			rc, ec := ca.session.SubmitRoomConfigurationForm(roomID, form)

			go func() {
				select {
				case <-rc:
					v.createReservedRoom(ca, roomID, onError)
				case errorResponse := <-ec:
					ca.log.WithError(errorResponse.Error()).Error("An error occurred when submitting the configuration form")
				case <-time.After(timeoutThreshold):
					onError(errCreateRoomTimeout)
				}
			}()
		case <-time.After(timeoutThreshold):
			onError(errCreateRoomTimeout)
		}
	}()
}

// joinRoom MUST be called from the UI thread
func (v *mucCreateRoomView) joinRoom(ca *account, roomID jid.Bare, d roomViewDataProvider) {
	v.dialog.Destroy()
	v.u.joinRoomWithData(ca, roomID, d)
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
	password     string
	autoJoin     bool
	customConfig bool
}

// passwordProvider implements the "roomViewDataProvider" interface
func (crd *mucCreateRoomData) passwordProvider() string {
	return crd.password
}

// backToPreviousStep implements the "roomViewDataProvider" interface
func (crd *mucCreateRoomData) backToPreviousStep() func() {
	return nil
}

// notifyError implements the "roomViewDataProvider" interface
func (crd *mucCreateRoomData) notifyError(err string) {
	// TODO: we need to check the current scenario to show the error notification.
	// 	1. Do we are in the create instant room scenario?
	// 	2. Do we are in the create a configured room scenario?
}

// doWhenNoErrors implements the "roomViewDataProvider" interface
func (crd *mucCreateRoomData) doWhenNoErrors() {
	// TODO: Close the current windows depending on the current scenario.
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
