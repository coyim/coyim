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
		c.log.Warn("Room operation succeed but no success callback was given")
		return
	}

	c.onSuccess()
}

func (c *roomOpController) error(err error) {
	log := c.log.WithError(err)
	if c.onError == nil {
		log.Error("Room operation failed but no error callback was given")
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

func (a *account) newAccountRoomOpContext(op string, roomID jid.Bare, c *roomOpController) *accountRoomOpContext {
	ctx := &accountRoomOpContext{
		op:         op,
		roomID:     roomID,
		account:    a,
		controller: c,
	}

	ctx.initLog()

	return ctx
}

func (ctx *accountRoomOpContext) initLog() {
	fields := ctx.getLogCommonFields()
	ctx.log = ctx.account.log.WithFields(fields)

	// We need the room view because if some error happens
	// during this operation we might want to log it using the room's logger
	room, exists := ctx.account.getRoomView(ctx.roomID)
	if !exists {
		// We don't call "stopWithError" here because we haven't even
		// started this context's operation
		ctx.controller.error(ctx.newInvalidRoomError())
		return
	}

	ctx.log = room.log.WithFields(fields)
}

func (ctx *accountRoomOpContext) getLogCommonFields() log.Fields {
	return log.Fields{
		"room":      ctx.roomID,
		"operation": ctx.op,
		"who":       "accountRoomOpContext",
	}
}

// doOperation will block until the controller finishes
func (ctx *accountRoomOpContext) doOperation() {
	ctx.successChannel = make(chan bool, 1)
	ctx.errorChannel = make(chan error, 1)
	ctx.cancelChannel = make(chan bool, 1)

	go ctx.waitUntilItFinish()
	ctx.controller.request(ctx.successChannel, ctx.errorChannel)
}

func (ctx *accountRoomOpContext) waitUntilItFinish() {
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
