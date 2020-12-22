package session

import (
	"errors"
	"fmt"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
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

// DestroyRoom sends a destroy room request, it receives:
//	- roomID Identifier of the room to be destroyed
//	- reason The reason why the room is being destroyed
//	- altRoomID The alternative room's jid where the discussions could continue
//	- password The alternative room password
// This method returns two read-only channels, one for the result and another one for any error
// that can happen during the room destruction process
func (s *session) DestroyRoom(roomID jid.Bare, reason string, altRoomID jid.Bare, password string) (<-chan bool, <-chan error) {
	rc := make(chan bool)
	ec := make(chan error)

	go s.muc.destroyRoom(roomID, newRoomDestroyQuery(reason, altRoomID, password), rc, ec)

	return rc, ec
}

func (m *mucManager) destroyRoom(roomID jid.Bare, q data.MUCRoomDestroyQuery, rc chan bool, ec chan error) {
	reply, _, err := m.conn().SendIQ(roomID.String(), "set", q)
	if err != nil {
		m.log.WithFields(log.Fields{
			"room":  roomID,
			"where": "destroyRoom",
		}).WithError(err).Error("Invalid destroy room information query response")
		ec <- err
		return
	}

	stanza := <-reply

	ciq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		ec <- ErrDestroyRoomInvalidIQResponse
		return
	}

	switch ciq.Type {
	case "result":
		m.deleteRoomFromManager(roomID)
		rc <- true
		return
	case "error":
		ec <- handleDestroyRoomIQError(ciq.Error)
		return
	default:
		ec <- ErrDestroyRoomNoResult
	}
}

// newRoomDestroyQuery returns a new query instance to be used as part of the destroy room process.
// This function receives:
//	- reason The reason of why the room is being destroyed
//	- altRoomID The alternative room identifier
//	- password The password to join the alternative room
func newRoomDestroyQuery(reason string, altRoomID jid.Bare, password string) data.MUCRoomDestroyQuery {
	return data.MUCRoomDestroyQuery{
		Destroy: data.MUCRoomDestroy{
			Reason:   reason,
			Jid:      fmt.Sprintf("%s", altRoomID),
			Password: password,
		},
	}
}

func handleDestroyRoomIQError(err data.StanzaError) error {
	switch {
	case err.Type == "cancel":
		return ErrDestroyRoomDoesntExist
	case err.MUCForbidden != nil:
		return ErrDestroyRoomForbidden
	default:
		return ErrDestroyRoomUnknown
	}
}
