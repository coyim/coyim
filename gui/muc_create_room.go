package gui

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/coylog"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucCreateRoomView struct {
	u *gtkUI

	autoJoin bool
	cancel   chan bool

	window    gtki.Window `gtk-widget:"createRoomWindow"`
	container gtki.Box    `gtk-widget:"content"`

	form    *mucCreateRoomViewForm
	success *mucCreateRoomViewSuccess

	showCreateForm  func()
	showSuccessView func(*account, jid.Bare)

	onAutoJoin *withCallbacks
	onDestroy  *withCallbacks

	sync.Mutex
}

func (u *gtkUI) newmucCreateRoomView() *mucCreateRoomView {
	v := &mucCreateRoomView{
		u:               u,
		showCreateForm:  func() {},
		showSuccessView: func(*account, jid.Bare) {},
		onAutoJoin:      newWithCallbacks(),
		onDestroy:       newWithCallbacks(),
	}

	v.initBuilder()

	v.form = v.initCreateRoomForm()
	v.showCreateForm = func() {
		v.form.showCreateForm(v)
	}

	v.success = v.initCreateRoomSuccess()
	v.showSuccessView = func(a *account, roomID jid.Bare) {
		v.success.showSuccessView(v, a, roomID)
	}

	u.connectShortcutsChildWindow(v.window)

	return v
}

func (v *mucCreateRoomView) initBuilder() {
	builder := newBuilder("MUCCreateRoomDialog")
	panicOnDevError(builder.bindObjects(v))

	builder.ConnectSignals(map[string]interface{}{
		"on_close_window": v.onCloseWindow,
	})
}

func (v *mucCreateRoomView) onCancel() {
	if v.cancel != nil {
		v.cancel <- true
		v.cancel = nil
	}

	v.window.Destroy()
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
	result := a.session.CreateRoom(roomID)
	go func() {
		err := <-result
		if err != nil {
			onError(err)
			return
		}
		onSuccess()
	}()
}

func (v *mucCreateRoomView) log(ca *account, roomID jid.Bare) coylog.Logger {
	var l coylog.Logger
	if ca != nil {
		l = ca.log
	} else {
		l = v.u.log
	}

	if roomID != nil {
		l.WithField("room", roomID)
	}

	l.WithField("who", "mucCreateRoomView")

	return l
}

func (v *mucCreateRoomView) createRoomIfDoesntExist(ca *account, roomID jid.Bare, errors chan error) {
	sc := make(chan bool)
	er := make(chan error)

	v.cancel = make(chan bool, 1)

	go func() {
		v.checkIfRoomExists(ca, roomID, sc, er)
		select {
		case <-sc:
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

func (v *mucCreateRoomView) onCreateRoomFinished(ca *account, roomID jid.Bare) {
	if !v.autoJoin {
		doInUIThread(func() {
			v.showSuccessView(ca, roomID)
			v.window.ShowAll()
		})
		return
	}

	v.joinRoom(ca, roomID)
}

func (v *mucCreateRoomView) joinRoom(ca *account, roomID jid.Bare) {
	doInUIThread(func() {
		v.window.Destroy()
		v.u.joinRoom(ca, roomID, nil)
	})
}

func (v *mucCreateRoomView) updateAutoJoinValue(f bool) {
	if v.autoJoin == f {
		return
	}

	v.Lock()
	defer v.Unlock()

	v.autoJoin = f
	v.onAutoJoin.invokeAll()
}

func (v *mucCreateRoomView) show() {
	v.showCreateForm()
	v.window.ShowAll()
}

func (u *gtkUI) mucCreateChatRoom() {
	view := u.newmucCreateRoomView()
	view.show()
}
