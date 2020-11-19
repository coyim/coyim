package gui

import (
	"fmt"

	"github.com/coyim/coyim/coylog"
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (a *account) getRoomView(roomID jid.Bare) (*roomView, bool) {
	a.mucRoomsLock.RLock()
	defer a.mucRoomsLock.RUnlock()

	v, ok := a.mucRooms[roomID.String()]
	if !ok {
		a.log.WithField("room", roomID).Debug("getRoomView(): trying to get a not connected room")
	}

	return v, ok
}

func (a *account) addRoomView(v *roomView) {
	a.mucRoomsLock.Lock()
	defer a.mucRoomsLock.Unlock()

	a.mucRooms[v.roomID().String()] = v
}

func (a *account) removeRoomView(roomID jid.Bare) {
	a.mucRoomsLock.Lock()
	defer a.mucRoomsLock.Unlock()

	_, exists := a.mucRooms[roomID.String()]
	if !exists {
		return
	}

	delete(a.mucRooms, roomID.String())
}

func (a *account) newRoomModel(roomID jid.Bare) *muc.Room {
	return a.session.NewRoom(roomID)
}

type roomOpCallback func() (<-chan bool, <-chan error, func())

type roomOpController struct {
	callback  roomOpCallback
	onSuccess func()
	onError   func(error)
	log       coylog.Logger
}

func (a *account) newRoomOpController(op string, cb roomOpCallback, onSuccess func(), onError func(error)) *roomOpController {
	return &roomOpController{
		callback:  cb,
		onSuccess: onSuccess,
		onError:   onError,
		log:       a.log.WithField("operation", op),
	}
}

func (c *roomOpController) request(sch chan bool, ech chan error) {
	ok, anyError, _ := c.callback()
	select {
	case <-ok:
		sch <- true
	case err := <-anyError:
		ech <- err
	}
}

func (c *roomOpController) success() {
	if c.onSuccess == nil {
		return
	}

	c.onSuccess()
}

func (c *roomOpController) error(err error) {
	c.log.WithError(err).Error("An error occurred while performing the operation in the room")

	if c.onError == nil {
		return
	}

	c.onError(err)
}

type accountRoomOpContext struct {
	op         string
	roomID     jid.Bare
	account    *account
	controller *roomOpController

	successChannel chan bool
	errorChannel   chan error
	cancelChannel  chan bool

	log coylog.Logger
}

func (a *account) newAccountRoomOpContext(op string, roomID jid.Bare, callback roomOpCallback, onSuccess func(), onError func(error)) *accountRoomOpContext {
	ctx := &accountRoomOpContext{
		op:         op,
		roomID:     roomID,
		account:    a,
		controller: a.newRoomOpController(op, callback, onSuccess, onError),
	}

	ctx.log = a.log.WithFields(log.Fields{
		"room":      ctx.roomID,
		"operation": ctx.op,
		"who":       "accountRoomOpContext",
	})

	return ctx
}

// doOperation will block until the controller finishes
func (ctx *accountRoomOpContext) doOperation() {
	ctx.successChannel = make(chan bool)
	ctx.errorChannel = make(chan error)
	ctx.cancelChannel = make(chan bool)

	go ctx.waitUntilFinished()
	ctx.controller.request(ctx.successChannel, ctx.errorChannel)
}

func (ctx *accountRoomOpContext) waitUntilFinished() {
	select {
	case <-ctx.successChannel:
		ctx.controller.success()
	case err := <-ctx.errorChannel:
		ctx.stopWithError(err)
	case <-ctx.cancelChannel:
	}
}

func (ctx *accountRoomOpContext) stopWithError(err error) {
	ctx.controller.error(err)
}

func (ctx *accountRoomOpContext) cancelOperation() {
	if ctx.cancelChannel == nil {
		return
	}

	ctx.log.Warn("A room operation was canceled, but can still occur")
	ctx.cancelChannel <- true
}

func (ctx *accountRoomOpContext) newInvalidRoomError() error {
	return fmt.Errorf("trying to %s a not available room \"%s\" for the account \"%s\"", ctx.op, ctx.roomID.String(), ctx.account.Account())
}
