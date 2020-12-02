package session

import (
	"errors"

	"github.com/coyim/coyim/coylog"

	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) DestroyRoom(roomID jid.Bare, reason string, alternativeRoomID jid.Bare, password string) (<-chan bool, <-chan error) {
	return s.muc.destroyRoom(roomID, reason, alternativeRoomID, password)
}

type destroyRoomRequest struct {
	roomID                  jid.Bare
	reason                  string
	alternativeRoomID       jid.Bare
	alternativeRoomPassword string
	onDestroyFinish         func(jid.Bare)

	resultChannel chan bool
	errorChannel  chan error

	conn xi.Conn
	log  coylog.Logger
}

func (m *mucManager) newDestroyRoomRequest(roomID jid.Bare, reason string, alternativeRoomID jid.Bare, password string, onDestroyFinish func(jid.Bare)) *destroyRoomRequest {
	return &destroyRoomRequest{
		roomID:                  roomID,
		conn:                    m.conn(),
		reason:                  reason,
		alternativeRoomID:       alternativeRoomID,
		alternativeRoomPassword: password,
		resultChannel:           make(chan bool),
		errorChannel:            make(chan error),
		onDestroyFinish:         onDestroyFinish,
		log:                     m.log,
	}
}

func (m *mucManager) destroyRoom(roomID jid.Bare, reason string, alternativeRoomID jid.Bare, password string) (<-chan bool, <-chan error) {
	m.destroyLock.Lock()
	defer m.destroyLock.Unlock()

	dr, ok := m.destroyRequests[roomID.String()]
	if !ok {
		dr = m.newDestroyRoomRequest(roomID, reason, alternativeRoomID, password, m.onDestroyRoomFinish)
		m.destroyRequests[roomID.String()] = dr

		go dr.sendDestroyRoomRequest()
	}

	return dr.resultChannel, dr.errorChannel
}

func (m *mucManager) onDestroyRoomFinish(roomID jid.Bare) {
	m.destroyLock.Lock()
	delete(m.destroyRequests, roomID.String())
	m.destroyLock.Unlock()
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

func (dr *destroyRoomRequest) sendDestroyRoomRequest() {
	defer dr.onDestroyFinish(dr.roomID)

	reply, err := dr.sendIQRequest()
	if err != nil {
		dr.finishWithError(err)
		return
	}

	err = dr.handleIQResponse(<-reply)
	if err != nil {
		dr.finishWithError(err)
		return
	}

	dr.finish()
}

func (dr *destroyRoomRequest) newRoomDestroyQuery() data.MUCRoomDestroyQuery {
	return data.MUCRoomDestroyQuery{
		Destroy: dr.newRoomDestroyData(),
	}
}

func (dr *destroyRoomRequest) newRoomDestroyData() data.MUCRoomDestroy {
	return data.MUCRoomDestroy{
		Reason:   dr.reason,
		Jid:      dr.alternativeRoomIDValue(),
		Password: dr.alternativeRoomPassword,
	}
}

func (dr *destroyRoomRequest) sendIQRequest() (<-chan data.Stanza, error) {
	q := dr.newRoomDestroyQuery()
	reply, cookie, err := dr.conn.SendIQ(dr.roomID.String(), "set", q)

	dr.log.WithField("cookie", cookie).Info("Sending an Information Query to destroy the room")

	return reply, err
}

func (dr *destroyRoomRequest) handleIQResponse(s data.Stanza) error {
	ciq, ok := s.Value.(*data.ClientIQ)
	if !ok {
		return ErrDestroyRoomInvalidIQResponse
	}

	switch ciq.Type {
	case "result":
		return nil
	case "error":
		return dr.handleIQError(ciq.Error)
	default:
		return ErrDestroyRoomNoResult
	}
}

func (dr *destroyRoomRequest) handleIQError(err data.StanzaError) error {
	switch {
	case err.MUCForbidden != nil:
		return ErrDestroyRoomForbidden
	default:
		return ErrDestroyRoomUnknown
	}
}

func (dr *destroyRoomRequest) finish() {
	dr.resultChannel <- true
}

func (dr *destroyRoomRequest) finishWithError(err error) {
	dr.log.WithError(err).Error("An error ocurred trying to destroy the room")
	dr.errorChannel <- err
}

func (dr *destroyRoomRequest) alternativeRoomIDValue() string {
	if dr.alternativeRoomID != nil {
		return dr.alternativeRoomID.String()
	}
	return ""
}
