package session

import (
	"errors"

	"github.com/coyim/coyim/coylog"
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

var (
	// ErrDestroyRoomInvalidIQResponse represents an invalid IQ response error
	ErrDestroyRoomInvalidIQResponse = errors.New("invalid destroy room IQ response")
	// ErrDestroyRoomForbidden represents a forbidden destroy room error
	ErrDestroyRoomForbidden = errors.New("destroy room forbidden")
	// ErrDestroyRoomUnknown represents an unknown destroy room error
	ErrDestroyRoomUnknown = errors.New("destroy room unknown error")
	// ErrDestroyRoomNoResult represents a no result received IQ error
	ErrDestroyRoomNoResult = errors.New("destroy room no result error")
)

func (s *session) DestroyRoom(roomID, alternativeRoomID jid.Bare, reason string) (<-chan bool, <-chan error, func()) {
	ctx := s.newDestroyRoomContext(roomID, alternativeRoomID, reason)

	go ctx.destroyRoom()

	return ctx.resultChannel, ctx.errorChannel, ctx.endEarly
}

type destroyRoomContext struct {
	roomID            jid.Bare
	alternativeRoomID jid.Bare
	reason            string

	resultChannel chan bool
	errorChannel  chan error
	cancelChannel chan bool

	conn xi.Conn
	log  coylog.Logger
}

func (s *session) newDestroyRoomContext(roomID, alternativeRoomID jid.Bare, reason string) *destroyRoomContext {
	return &destroyRoomContext{
		roomID:            roomID,
		alternativeRoomID: alternativeRoomID,
		reason:            reason,
		resultChannel:     make(chan bool, 1),
		errorChannel:      make(chan error, 1),
		conn:              s.conn,
		log: s.log.WithFields(log.Fields{
			"room":            roomID,
			"alternativeRoom": alternativeRoomID,
			"context":         "destroy-room",
		}),
	}
}

func (ctx *destroyRoomContext) destroyRoom() {
	if ctx.cancelChannel != nil {
		ctx.finishWithCancel()
		return
	}

	ctx.cancelChannel = make(chan bool, 1)

	reply, err := ctx.sendIQRequest()
	if err != nil {
		ctx.finishWithError(err)
		return
	}

	select {
	case s := <-reply:
		ctx.onIQResponse(s)
	case <-ctx.cancelChannel:
		ctx.finishWithCancel()
	}
}

func (ctx *destroyRoomContext) newRoomDestroyQuery() data.MUCRoomDestroyQuery {
	return data.MUCRoomDestroyQuery{
		Destroy: ctx.newRoomDestroyData(),
	}
}

func (ctx *destroyRoomContext) newRoomDestroyData() data.MUCRoomDestroy {
	return data.MUCRoomDestroy{
		Jid:    ctx.getAlternativeRoomID(),
		Reason: ctx.reason,
	}
}

func (ctx *destroyRoomContext) sendIQRequest() (<-chan data.Stanza, error) {
	q := ctx.newRoomDestroyQuery()
	reply, cookie, err := ctx.conn.SendIQ(ctx.roomID.String(), "set", q)

	ctx.log.WithField("cookie", cookie).Info("Sending an Information Query to destroy the room")

	return reply, err
}

func (ctx *destroyRoomContext) onIQResponse(s data.Stanza) {
	err := ctx.handleIQResponse(s)
	if err != nil {
		ctx.finishWithError(err)
		return
	}

	ctx.finish()
}

func (ctx *destroyRoomContext) handleIQResponse(s data.Stanza) error {
	ciq, ok := s.Value.(*data.ClientIQ)
	if !ok {
		return ErrDestroyRoomInvalidIQResponse
	}

	switch ciq.Type {
	case "result":
		return nil
	case "error":
		return ctx.handleIQError(ciq.Error)
	default:
		return ErrDestroyRoomNoResult
	}
}

func (ctx *destroyRoomContext) handleIQError(err data.StanzaError) error {
	switch {
	case err.MUCForbidden != nil:
		return ErrDestroyRoomForbidden
	default:
		return ErrDestroyRoomUnknown
	}
}

func (ctx *destroyRoomContext) finish() {
	ctx.resultChannel <- true
	ctx.clean()
}

func (ctx *destroyRoomContext) finishWithError(err error) {
	ctx.log.WithError(err).Error("An error ocurred trying to destroy the room")
	ctx.errorChannel <- err
	ctx.clean()
}

func (ctx *destroyRoomContext) finishWithCancel() {
	ctx.log.Warn("The destroy room operation was canceled, but it could still happen")
	ctx.clean()
}

func (ctx *destroyRoomContext) endEarly() {
	if ctx.cancelChannel == nil {
		ctx.cancelChannel = make(chan bool, 1)
	}

	ctx.cancelChannel <- true
}

func (ctx *destroyRoomContext) clean() {
	close(ctx.resultChannel)
	close(ctx.errorChannel)
	close(ctx.cancelChannel)
}

func (ctx *destroyRoomContext) getAlternativeRoomID() string {
	if ctx.alternativeRoomID != nil {
		return ctx.alternativeRoomID.String()
	}
	return ""
}
