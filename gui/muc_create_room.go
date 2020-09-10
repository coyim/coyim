package gui

import (
	"errors"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoom struct {
	u       *gtkUI
	builder *builder
	log     coylog.Logger

	autoJoin bool

	window    gtki.Window `gtk-widget:"createRoomWindow"`
	container gtki.Box    `gtk-widget:"content"`

	form    *createMUCRoomForm
	success *createMUCRoomSuccess

	showCreateForm  func()
	showSuccessView func(*account, jid.Bare)
	onAutoJoinList  []func(bool)
	onDestroyList   []func()
	builderSignals  map[string]interface{}
}

func (u *gtkUI) newCreateMUCRoom() *createMUCRoom {
	v := &createMUCRoom{
		u:               u,
		log:             u.log,
		showCreateForm:  func() {},
		showSuccessView: func(*account, jid.Bare) {},
		onAutoJoinList:  []func(bool){},
		onDestroyList:   []func(){},
		builderSignals:  make(map[string]interface{}),
	}

	v.initUIBuilder()
	v.initForm()
	v.initSuccessView(v.joinRoom)
	v.initBuilderSignals()

	u.connectShortcutsChildWindow(v.window)

	return v
}

func (v *createMUCRoom) updateAutoJoinValue(newValue bool) {
	if v.autoJoin == newValue {
		return
	}

	v.autoJoin = newValue
	for _, cb := range v.onAutoJoinList {
		cb(v.autoJoin)
	}
}

func (v *createMUCRoom) onAutoJoin(f func(bool)) {
	v.onAutoJoinList = append(v.onAutoJoinList, f)
}

func (v *createMUCRoom) onDestroy(f func()) {
	v.onDestroyList = append(v.onDestroyList, f)
}

func (v *createMUCRoom) addBuilderSignals(signals map[string]interface{}) {
	for signal, callback := range signals {
		v.addBuilderSignal(signal, callback)
	}
}

func (v *createMUCRoom) addBuilderSignal(signal string, callback interface{}) {
	if _, ok := v.builderSignals[signal]; ok {
		v.log.WithField("signal", signal).Warn("Signal already registered")
		return
	}
	v.builderSignals[signal] = callback
}

func (v *createMUCRoom) initUIBuilder() {
	v.builder = newBuilder("MUCCreateRoomDialog")
	panicOnDevError(v.builder.bindObjects(v))
}

func (v *createMUCRoom) initBuilderSignals() {
	v.addBuilderSignal("on_close_window", v.onCloseWindow)
	v.builder.ConnectSignals(v.builderSignals)
}

func (v *createMUCRoom) onCancel() {
	v.window.Destroy()
}

func (v *createMUCRoom) onCloseWindow() {
	for _, cb := range v.onDestroyList {
		cb()
	}
}

var (
	errCreateRoomCheckIfExistsFails = errors.New("room exists failed")
	errCreateRoomAlreadyExists      = errors.New("room already exists")
	errCreateRoomFailed             = errors.New("couldn't create the room")
)

func (v *createMUCRoom) createRoomIfDoesntExist(ca *account, ident jid.Bare, successResult chan bool, errResult chan error) {
	erc, ec := ca.session.HasRoom(ident, nil)
	go func() {
		select {
		case err := <-ec:
			if err != nil {
				ca.log.WithError(err).Error("Error trying to validate if room exists")
				errResult <- errCreateRoomCheckIfExistsFails
			}

		case er := <-erc:
			if !er {
				ec := ca.session.CreateRoom(ident)
				go func() {
					err := <-ec
					if err != nil {
						ca.log.WithError(err).Error("Something went wrong while trying to create the room")
						errResult <- errCreateRoomFailed
						return
					}
					v.onCreateRoomFinished(ca, ident)
				}()
				return
			}
			errResult <- errCreateRoomAlreadyExists
		}
	}()
}

func (v *createMUCRoom) onCreateRoomFinished(ca *account, ident jid.Bare) {
	if !v.autoJoin {
		doInUIThread(func() {
			v.showSuccessView(ca, ident)
		})
		return
	}

	v.joinRoom(ca, ident)
}

func (v *createMUCRoom) joinRoom(ca *account, ident jid.Bare) {
	doInUIThread(v.destroy)
	go func() {
		rl := make(chan *muc.RoomListing)
		go ca.session.GetRoom(ident, rl)
		roomInfo := <-rl

		doInUIThread(func() {
			v.u.mucShowRoom(ca, ident, roomInfo, nil)
		})
	}()
}

func (v *createMUCRoom) destroy() {
	v.window.Destroy()
}

func (v *createMUCRoom) show() {
	v.showCreateForm()
	v.window.Show()
}

func (u *gtkUI) mucCreateChatRoom() {
	view := u.newCreateMUCRoom()
	view.show()
}
