package gui

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoom struct {
	u   *gtkUI
	log coylog.Logger

	autoJoin bool
	cancel   chan bool

	window    gtki.Window `gtk-widget:"createRoomWindow"`
	container gtki.Box    `gtk-widget:"content"`

	form    *createMUCRoomForm
	success *createMUCRoomSuccess

	showCreateForm   func()
	showSuccessView  func(*account, jid.Bare)
	onAutoJoinList   []func(bool)
	onDestroyList    []func()
	onAutoJoinLocker sync.RWMutex
	onDestroyLocker  sync.RWMutex
}

func (u *gtkUI) newCreateMUCRoom() *createMUCRoom {
	v := &createMUCRoom{
		u:               u,
		log:             u.log,
		showCreateForm:  func() {},
		showSuccessView: func(*account, jid.Bare) {},
		onAutoJoinList:  []func(bool){},
		onDestroyList:   []func(){},
	}

	v.initBuilder()

	v.form = v.newCreateRoomForm()
	v.success = v.newCreateRoomSuccess()

	u.connectShortcutsChildWindow(v.window)

	return v
}

func (v *createMUCRoom) initBuilder() {
	builder := newBuilder("MUCCreateRoomDialog")
	panicOnDevError(builder.bindObjects(v))

	builder.ConnectSignals(map[string]interface{}{
		"on_close_window": v.onCloseWindow,
	})
}

func (v *createMUCRoom) onDestroy(f func()) {
	v.onDestroyLocker.Lock()
	defer v.onDestroyLocker.Unlock()

	v.onDestroyList = append(v.onDestroyList, f)
}

func (v *createMUCRoom) onCancel() {
	if v.cancel != nil {
		v.cancel <- true
	}

	v.window.Destroy()
}

func (v *createMUCRoom) onCloseWindow() {
	v.onDestroyLocker.RLock()
	defer v.onDestroyLocker.RUnlock()

	for _, cb := range v.onDestroyList {
		cb()
	}
}

var (
	errCreateRoomCheckIfExistsFails = errors.New("room exists failed")
	errCreateRoomAlreadyExists      = errors.New("room already exists")
	errCreateRoomFailed             = errors.New("couldn't create the room")
)

func (v *createMUCRoom) checkIfRoomExists(ca *account, ident jid.Bare, result chan bool, errors chan error) {
	rc, ec := ca.session.HasRoom(ident, nil)
	go func() {
		select {
		case err := <-ec:
			ca.log.WithError(err).Error("Error trying to validate if room exists")
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

func (a *account) createRoom(ident jid.Bare, onSuccess func(), onError func(error)) {
	result := a.session.CreateRoom(ident)
	go func() {
		err := <-result
		if err != nil {
			onError(err)
			return
		}
		onSuccess()
	}()
}

func (v *createMUCRoom) createRoomIfDoesntExist(ca *account, ident jid.Bare, errors chan error) {
	sc := make(chan bool)
	er := make(chan error)

	v.cancel = make(chan bool, 1)

	go func() {
		v.checkIfRoomExists(ca, ident, sc, er)
		select {
		case <-sc:
			ca.createRoom(ident, func() {
				v.onCreateRoomFinished(ca, ident)
			}, func(err error) {
				ca.log.WithError(err).Error("Something went wrong while trying to create the room")
				errors <- errCreateRoomFailed
			})
		case err := <-er:
			errors <- err
		case <-v.cancel:
		}
	}()
}

func (v *createMUCRoom) onCreateRoomFinished(ca *account, ident jid.Bare) {
	if !v.autoJoin {
		doInUIThread(func() {
			v.showSuccessView(ca, ident)
			v.window.ShowAll()
		})
		return
	}

	v.joinRoom(ca, ident)
}

func (v *createMUCRoom) joinRoom(ca *account, ident jid.Bare) {
	doInUIThread(func() {
		v.destroy()
		v.u.joinMultiUserChat(ca, ident, nil)
	})
}

func (v *createMUCRoom) updateAutoJoinValue(newValue bool) {
	if v.autoJoin == newValue {
		return
	}

	v.onAutoJoinLocker.RLock()
	defer v.onAutoJoinLocker.RUnlock()

	v.autoJoin = newValue
	for _, cb := range v.onAutoJoinList {
		cb(v.autoJoin)
	}
}

func (v *createMUCRoom) onAutoJoin(f func(bool)) {
	v.onAutoJoinLocker.Lock()
	defer v.onAutoJoinLocker.Unlock()

	v.onAutoJoinList = append(v.onAutoJoinList, f)
}

func (v *createMUCRoom) destroy() {
	v.window.Destroy()
}

func (v *createMUCRoom) show() {
	v.showCreateForm()
	v.window.ShowAll()
}

func (u *gtkUI) mucCreateChatRoom() {
	view := u.newCreateMUCRoom()
	view.show()
}
