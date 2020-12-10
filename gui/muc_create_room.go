package gui

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
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

	showCreateForm  func()
	showSuccessView func(*account, jid.Bare)

	onCreateOptionChange *callbacksSet
	onDestroy            *callbacksSet

	sync.Mutex
}

func newCreateMUCRoomView(u *gtkUI) *mucCreateRoomView {
	v := &mucCreateRoomView{
		u:                    u,
		showCreateForm:       func() {},
		showSuccessView:      func(*account, jid.Bare) {},
		onCreateOptionChange: newCallbacksSet(),
		onDestroy:            newCallbacksSet(),
	}

	v.initBuilder()
	v.initChildViews()

	return v
}

func (v *mucCreateRoomView) initBuilder() {
	builder := newBuilder("MUCCreateRoomDialog")
	panicOnDevError(builder.bindObjects(v))

	builder.ConnectSignals(map[string]interface{}{
		"on_close_window": v.onCloseWindow,
	})
}

func (v *mucCreateRoomView) initChildViews() {
	v.form = v.initCreateRoomForm()
	v.showCreateForm = func() {
		v.form.showCreateForm(v)
	}

	v.success = v.initCreateRoomSuccess()
	v.showSuccessView = func(ca *account, roomID jid.Bare) {
		v.success.showSuccessView(v, ca, roomID)
	}
}

func (v *mucCreateRoomView) onCancel() {
	if v.cancel != nil {
		v.cancel <- true
		v.cancel = nil
	}

	v.dialog.Destroy()
}

func (v *mucCreateRoomView) onCloseWindow() {
	v.onDestroy.invokeAll()
}

var (
	errCreateRoomCheckIfExistsFails = errors.New("room exists failed")
	errCreateRoomAlreadyExists      = errors.New("room already exists")
	errCreateRoomFailed             = errors.New("couldn't create the room")
)

func (v *mucCreateRoomView) checkIfRoomExists(ca *account, roomID jid.Bare, result chan bool, errors chan error) {
	rc, ec := ca.session.HasRoom(roomID, nil)
	go func() {
		select {
		case err := <-ec:
			v.log(ca, roomID).WithError(err).Error("Error trying to validate if room exists")
			errors <- errCreateRoomCheckIfExistsFails
		case exists := <-rc:
			if exists {
				errors <- errCreateRoomAlreadyExists
				return
			}
			result <- true
		case <-v.cancel:
		}
	}()
}

func (a *account) createRoom(roomID jid.Bare, onSuccess func(), onError func(error)) {
	result := a.session.CreateInstantRoom(roomID)
	go func() {
		err := <-result
		if err != nil {
			onError(err)
			return
		}
		onSuccess()
	}()
}

func (a *account) reserveRoom(roomID jid.Bare, onSuccess func(*muc.RoomConfigForm), onError func(error)) {
	fc, ec := a.session.CreateReservedRoom(roomID)
	go func() {
		select {
		case err := <-ec:
			if err != nil {
				onError(err)
				return
			}
		case form := <-fc:
			onSuccess(form)
		}
	}()
}

func (v *mucCreateRoomView) log(ca *account, roomID jid.Bare) coylog.Logger {
	l := v.u.log
	if ca != nil {
		l = ca.log
	}

	if roomID != nil {
		l.WithField("room", roomID)
	}

	l.WithField("where", "mucCreateRoomView")

	return l
}

func (v *mucCreateRoomView) createRoom(ca *account, roomID jid.Bare, errors chan error) {
	sc := make(chan bool)
	er := make(chan error)

	v.cancel = make(chan bool, 1)

	go func() {
		v.checkIfRoomExists(ca, roomID, sc, er)
		select {
		case <-sc:
			if v.configureRoom {
				ca.reserveRoom(roomID, v.onReserveRoomFinished, func(err error) {
					v.log(ca, roomID).WithError(err).Error("Something went wrong while trying to create the room")
					errors <- errCreateRoomFailed
				})
				return
			}

			ca.createRoom(roomID, func() {
				v.onCreateRoomFinished(ca, roomID)
			}, func(err error) {
				v.log(ca, roomID).WithError(err).Error("Something went wrong while trying to create the room")
				errors <- errCreateRoomFailed
			})
		case err := <-er:
			errors <- err
		case <-v.cancel:
		}
	}()
}

func (v *mucCreateRoomView) onReserveRoomFinished(cf *muc.RoomConfigForm) {
	doInUIThread(func() {
		v.dialog.Destroy()
		rca := newRoomConfigAssistant(cf)
		rca.show()
	})
}

func (v *mucCreateRoomView) onCreateRoomFinished(ca *account, roomID jid.Bare) {
	if v.autoJoin {
		doInUIThread(func() {
			v.joinRoom(ca, roomID)
		})
		return
	}

	doInUIThread(func() {
		v.showSuccessView(ca, roomID)
		v.dialog.ShowAll()
	})
}

// joinRoom MUST be called from the UI thread
func (v *mucCreateRoomView) joinRoom(ca *account, roomID jid.Bare) {
	v.dialog.Destroy()
	v.u.joinRoom(ca, roomID, nil)
}

func (v *mucCreateRoomView) updateAutoJoinValue(f bool) {
	v.updateCreateOption("autoJoin", f)
}

func (v *mucCreateRoomView) updateConfigureRoomValue(f bool) {
	v.updateCreateOption("configRoom", f)
}

func (v *mucCreateRoomView) updateCreateOption(o string, f bool) {
	v.Lock()
	defer v.Unlock()

	oldValue := false

	switch o {
	case "autoJoin":
		oldValue = v.autoJoin
		v.autoJoin = f
	case "configRoom":
		oldValue = v.configureRoom
		v.configureRoom = f
	}

	if oldValue != f {
		v.onCreateOptionChange.invokeAll()
	}
}

func (u *gtkUI) mucCreateChatRoom() {
	view := newCreateMUCRoomView(u)

	u.connectShortcutsChildWindow(view.dialog)

	view.showCreateForm()

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
