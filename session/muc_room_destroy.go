package session

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/coylog"
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) DestroyRoom(roomID jid.Bare, reason string, alternativeRoomID jid.Bare, password string) (<-chan bool, <-chan error) {
	return s.muc.destroyRoom(roomID, reason, alternativeRoomID, password)
}

type destroyRequest struct {
	conn           xi.Conn
	roomID         jid.Bare
	destroyContext []*destroyRoomContext
	lock           sync.Mutex
	log            coylog.Logger
}

func (dr *destroyRequest) resolveNewRequest(reason string, alternativeRoomID jid.Bare, password string, resultChannel chan bool, errorChannel chan error) {
	dr.lock.Lock()
	defer dr.lock.Unlock()

	ctx := dr.newDestroyRoomContext(reason, alternativeRoomID, password)

	go ctx.destroyRoom()

	select {
	case <-ctx.resultChannel:
		resultChannel <- true
	case err := <-ctx.errorChannel:
		errorChannel <- err
	}
}

func (m *mucManager) destroyRoom(roomID jid.Bare, reason string, alternativeRoomID jid.Bare, password string) (<-chan bool, <-chan error) {
	dr := m.requestForRoomID(roomID)

	m.destroyLock.Lock()
	defer m.destroyLock.Unlock()

	rc := make(chan bool)
	ec := make(chan error)

	go dr.resolveNewRequest(reason, alternativeRoomID, password, rc, ec)

	return rc, ec
}

func (m *mucManager) requestForRoomID(roomID jid.Bare) *destroyRequest {
	m.destroyLock.RLock()
	defer m.destroyLock.RUnlock()

	for _, dr := range m.destroyRequests {
		if dr.roomID.String() == roomID.String() {
			return dr
		}
	}

	dr := &destroyRequest{
		conn:   m.conn(),
		roomID: roomID,
		log: m.log.WithFields(log.Fields{
			"room":    roomID,
			"request": "destroy-room",
		}),
	}

	m.destroyRequests = append(m.destroyRequests, dr)

	return dr
}

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

type destroyRoomContext struct {
	roomID                  jid.Bare
	reason                  string
	alternativeRoomID       jid.Bare
	alternativeRoomPassword string

	resultChannel chan bool
	errorChannel  chan error

	conn xi.Conn
	log  coylog.Logger
}

func (dr *destroyRequest) newDestroyRoomContext(reason string, alternativeRoomID jid.Bare, password string) *destroyRoomContext {
	return &destroyRoomContext{
		roomID:                  dr.roomID,
		conn:                    dr.conn,
		reason:                  reason,
		alternativeRoomID:       alternativeRoomID,
		alternativeRoomPassword: password,
		resultChannel:           make(chan bool),
		errorChannel:            make(chan error),
		log:                     dr.log,
	}
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

	ctx.finish()
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
	q := ctx.newRoomDestroyQuery()
	reply, cookie, err := ctx.conn.SendIQ(ctx.roomID.String(), "set", q)

	ctx.log.WithField("cookie", cookie).Info("Sending an Information Query to destroy the room")

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
	switch {
	case err.MUCForbidden != nil:
		return ErrDestroyRoomForbidden
	default:
		return ErrDestroyRoomUnknown
	}
}

func (ctx *destroyRoomContext) finish() {
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
