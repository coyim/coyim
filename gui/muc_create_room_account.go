package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

// createInstantRoom IS SAFE to be called from the UI thread
func (a *account) createInstantRoom(roomID jid.Bare, onSuccess func(), onError func(error)) {
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

// createReservedRoom IS SAFE to be called from the UI thread
func (a *account) createReservedRoom(roomID jid.Bare, onSuccess func(jid.Bare, *muc.RoomConfigForm), onError func(error)) {
	fc, ec := a.session.CreateReservedRoom(roomID)
	go func() {
		select {
		case err := <-ec:
			if err != nil {
				onError(err)
				return
			}
		case form := <-fc:
			onSuccess(roomID, form)
		}
	}()
}

func (v *mucCreateRoomView) createReservedRoom(ca *account, roomID jid.Bare, onError func(error)) {
	ca.createReservedRoom(roomID, func(roomID jid.Bare, cf *muc.RoomConfigForm) {
		doInUIThread(func() {
			v.onReserveRoomFinished(ca, roomID, cf)
		})
	}, func(err error) {
		v.log(ca, roomID).WithError(err).Error("Something went wrong when trying to reserve the room")
		onError(errCreateRoomFailed)
	})
}

func (v *mucCreateRoomView) createInstantRoom(ca *account, roomID jid.Bare, onError func(error)) {
	ca.createInstantRoom(roomID, func() {
		v.onCreateRoomFinished(ca, roomID, v.autoJoin, func() {
			v.showSuccessView(ca, roomID)
			v.dialog.ShowAll()
		})
	}, func(err error) {
		v.log(ca, roomID).WithError(err).Error("Something went wrong when trying to create the instant room")
		onError(errCreateRoomFailed)
	})
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

// onReserveRoomFinished MUST be called from the UI thread
func (v *mucCreateRoomView) onReserveRoomFinished(ca *account, roomID jid.Bare, cf *muc.RoomConfigForm) {
	rca := v.u.newRoomConfigAssistant(ca, roomID, cf, v.autoJoin, func(ca *account, roomID jid.Bare, autoJoin bool) {
		v.onCreateRoomFinished(ca, roomID, autoJoin, func() {
			v.u.mucShowCreateRoomSuccess(ca, roomID, &mucCreateRoomData{
				ca:           ca,
				where:        roomID.Host(),
				customConfig: true,
			})
		})
	}, func() {
		doInUIThread(func() {
			v.u.mucShowCreateRoomForm(&mucCreateRoomData{
				ca:           ca,
				roomName:     roomID.Local(),
				where:        roomID.Host(),
				autoJoin:     v.autoJoin,
				customConfig: true,
			})
		})
	})
	rca.show()
	v.destroy()
}

// onCreateRoomFinished MUST NOT be called from the UI thread
func (v *mucCreateRoomView) onCreateRoomFinished(ca *account, roomID jid.Bare, autoJoin bool, onNoAutoJoin func()) {
	if autoJoin {
		doInUIThread(func() {
			v.joinRoom(ca, roomID)
		})
		return
	}

	if onNoAutoJoin != nil {
		doInUIThread(onNoAutoJoin)
	}
}
