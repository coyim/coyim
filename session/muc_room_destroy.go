package session

import (
	"errors"

	"github.com/coyim/coyim/coylog"

	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

var (
	// ErrDestroyRoomInvalidIQResponse represents an invalid IQ response error
	ErrDestroyRoomInvalidIQResponse = errors.New("invalid destroy room IQ response")
	// ErrDestroyRoomForbidden represents a forbidden destroy room error
	ErrDestroyRoomForbidden = errors.New("destroy room forbidden")
	// ErrDestroyRoomDoesntExist represents an unknown destroy room error
	ErrDestroyRoomDoesntExist = errors.New("room doesn't exist")
	// ErrDestroyRoomUnknown represents an unknown destroy room error
	ErrDestroyRoomUnknown = errors.New("destroy room unknown error")
	// ErrDestroyRoomNoResult represents a no result received IQ error
	ErrDestroyRoomNoResult = errors.New("destroy room no result error")
)

type destroyRoomContext struct {
	roomID                  jid.Bare
	reason                  string
	alternativeRoomID       jid.Bare
	alternativeRoomPassword string
	onSuccess               func(roomID jid.Bare) bool

	resultChannel chan bool
	errorChannel  chan error

	conn xi.Conn
	log  coylog.Logger
}

func (s *session) newDestroyRoomContext(roomID jid.Bare, reason string, alternativeRoomID jid.Bare, password string, onSuccess func(jid.Bare) bool) *destroyRoomContext {
	return &destroyRoomContext{
		roomID:                  roomID,
		reason:                  reason,
		alternativeRoomID:       alternativeRoomID,
		alternativeRoomPassword: password,
		onSuccess:               onSuccess,
		resultChannel:           make(chan bool),
		errorChannel:            make(chan error),
		conn:                    s.conn,
		log:                     s.log.WithField("where", "destroyRoomContext"),
	}
}

func (s *session) DestroyRoom(roomID jid.Bare, reason string, alternativeRoomID jid.Bare, password string) (<-chan bool, <-chan error) {
	ctx := s.newDestroyRoomContext(roomID, reason, alternativeRoomID, password, s.muc.roomManager.LeaveRoom)
	go ctx.destroyRoom()
	return ctx.resultChannel, ctx.errorChannel
}

func (ctx *destroyRoomContext) destroyRoom() {
	reply, err := ctx.sendIQRequest()
	if err != nil {
		ctx.finishWithError(err)
		return
	}

	err = ctx.handleIQResponse(<-reply)
	if err != nil {
		ctx.finishWithError(err)
		return
	}

	ctx.finishWithSuccess()
}

func (ctx *destroyRoomContext) newRoomDestroyQuery() data.MUCRoomDestroyQuery {
	return data.MUCRoomDestroyQuery{
		Destroy: ctx.newRoomDestroyData(),
	}
}

func (ctx *destroyRoomContext) newRoomDestroyData() data.MUCRoomDestroy {
	return data.MUCRoomDestroy{
		Reason:   ctx.reason,
		Jid:      ctx.alternativeRoomIDValue(),
		Password: ctx.alternativeRoomPassword,
	}
}

func (ctx *destroyRoomContext) sendIQRequest() (<-chan data.Stanza, error) {
	reply, _, err := ctx.conn.SendIQ(ctx.roomID.String(), "set", ctx.newRoomDestroyQuery())
	return reply, err
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
	if err.Type == "cancel" {
		return ErrDestroyRoomDoesntExist
	}

	switch {
	case err.MUCForbidden != nil:
		return ErrDestroyRoomForbidden
	default:
		return ErrDestroyRoomUnknown
	}
}

func (ctx *destroyRoomContext) finishWithSuccess() {
	ctx.onSuccess(ctx.roomID)
	ctx.resultChannel <- true
}

func (ctx *destroyRoomContext) finishWithError(err error) {
	ctx.log.WithError(err).Error("An error ocurred trying to destroy the room")
	ctx.errorChannel <- err
}

func (ctx *destroyRoomContext) alternativeRoomIDValue() string {
	if ctx.alternativeRoomID != nil {
		return ctx.alternativeRoomID.String()
	}
	return ""
}
