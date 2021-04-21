package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

// createInstantRoom IS SAFE to be called from the UI thread
func (a *account) createInstantRoom(roomID jid.Bare, onSuccess func(), onError func(error)) {
	rc, ec := a.session.CreateInstantRoom(roomID)
	go func() {
		select {
		case err := <-ec:
			onError(err)
		case <-rc:
			onSuccess()
		}
	}()
}

// createReservedRoom IS SAFE to be called from the UI thread
func (a *account) createReservedRoom(roomID jid.Bare, onSuccess func(jid.Bare, *muc.RoomConfigForm), onError func(error)) {
	fc, ec := a.session.CreateReservedRoom(roomID)
	go func() {
		select {
		case err := <-ec:
			onError(err)
		case form := <-fc:
			onSuccess(roomID, form)
		}
	}()
}

func (v *mucCreateRoomView) createReservedRoom(ca *account, roomID jid.Bare, onError func(error)) {
	onSuccess := func(roomID jid.Bare, cf *muc.RoomConfigForm) {
		doInUIThread(func() {
			v.onReserveRoomFinished(ca, roomID, cf)
		})
	}

	onErrorFinal := onError
	onError = func(err error) {
		v.log(ca, roomID).WithError(err).Error("Something went wrong when trying to reserve the room")
		onErrorFinal(errCreateRoomFailed)
	}

	ca.createReservedRoom(roomID, onSuccess, onError)
}

func (v *mucCreateRoomView) createInstantRoom(ca *account, roomID jid.Bare, onError func(error)) {
	d := newCreateRoomData()
	d.autoJoin = v.autoJoin

	onSuccess := func() {
		v.onCreateRoomFinished(ca, roomID, d, func() {
			v.showSuccessView(ca, roomID)
			v.dialog.ShowAll()
		})
	}

	onErrorFinal := onError
	onError = func(err error) {
		v.log(ca, roomID).WithError(err).Error("Something went wrong when trying to create the instant room")
		onErrorFinal(errCreateRoomFailed)
	}

	ca.createInstantRoom(roomID, onSuccess, onError)
}

// checkIfRoomExists IS SAFE to be called from the UI thread
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

func (v *mucCreateRoomView) createRoomDataBasedOnConfigForm(ca *account, roomID jid.Bare, cf *muc.RoomConfigForm) *mucCreateRoomData {
	crd := newCreateRoomData()

	crd.ca = ca
	crd.roomName = roomID.Local()
	crd.where = roomID.Host()
	crd.password = cf.Password
	crd.autoJoin = v.autoJoin
	crd.customConfig = true

	return crd
}

// onReserveRoomFinished MUST be called from the UI thread
func (v *mucCreateRoomView) onReserveRoomFinished(ca *account, roomID jid.Bare, cf *muc.RoomConfigForm) {
	v.destroy()

	createRoomData := v.createRoomDataBasedOnConfigForm(ca, roomID, cf)

	onSuccess := func(autoJoin bool) {
		createRoomData.autoJoin = autoJoin
		createRoomData.password = cf.Password
		v.onCreateRoomFinished(ca, roomID, createRoomData, func() {
			v.u.mucShowCreateRoomSuccess(ca, roomID, createRoomData)
		})
	}

	onCancel := func() {
		doInUIThread(func() {
			v.u.mucShowCreateRoomForm(createRoomData)
		})
	}

	rca := v.u.newRoomConfigAssistant(ca, roomID, cf, v.autoJoin, onSuccess, onCancel)
	rca.showAssistant()
}

// onCreateRoomFinished MUST NOT be called from the UI thread
func (v *mucCreateRoomView) onCreateRoomFinished(ca *account, roomID jid.Bare, createRoomData *mucCreateRoomData, onNoAutoJoin func()) {
	if createRoomData.autoJoin {
		doInUIThread(func() {
			v.joinRoom(ca, roomID, createRoomData)
		})
		return
	}

	if onNoAutoJoin != nil {
		doInUIThread(onNoAutoJoin)
	}
}
