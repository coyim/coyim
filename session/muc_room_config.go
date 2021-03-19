package session

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrRoomConfigSubmit represents an early error that happened during room configuration submit
	ErrRoomConfigSubmit = errors.New("invalid room configuration submit request")
	// ErrRoomConfigSubmitResponse represents an invalid response for a room configuration request
	ErrRoomConfigSubmitResponse = errors.New("invalid response for room configuration request")
	// ErrRoomConfigCancel represents an early error that happened during a room configuration cancel request
	ErrRoomConfigCancel = errors.New("invalid room configuration cancel request")
	// ErrRoomConfigCancelResponse represents an invalid response for a room configuration cancel request
	ErrRoomConfigCancelResponse = errors.New("invalid response for room configuration cancel request")
)

func (s *session) SubmitRoomConfigurationForm(roomID jid.Bare, form *muc.RoomConfigForm) (<-chan bool, <-chan error) {
	log := log.WithFields(log.Fields{
		"room":  roomID,
		"where": "SubmitRoomConfigurationForm",
	})

	sc := make(chan bool)
	ec := make(chan error)

	go func() {
		reply, _, err := s.conn.SendIQ(roomID.String(), "set", data.MUCRoomConfiguration{
			Form: form.GetFormData(),
		})

		if err != nil {
			log.WithError(err).Error("An error occurred when trying to send the information query to save the room configuration")
			ec <- ErrRoomConfigSubmit
			return
		}

		err = validateIqResponse(reply)
		if err != nil {
			log.WithError(ErrInformationQueryResponse).Error("An error occurred when trying to read the response from the room configuration request")
			ec <- ErrRoomConfigSubmitResponse
			return
		}

		sc <- true

	}()

	return sc, ec
}

func (s *session) CancelRoomConfiguration(roomID jid.Bare) <-chan error {
	log := log.WithFields(log.Fields{
		"room":  roomID,
		"where": "CancelRoomConfiguration",
	})

	ec := make(chan error)

	go func() {
		reply, _, err := s.conn.SendIQ(roomID.String(), "set", data.MUCRoomConfiguration{
			Form: &data.Form{
				Type: "cancel",
			},
		})

		if err != nil {
			log.WithError(err).Error("An error occurred when trying to send the request to rollback the room configuration")
			ec <- ErrRoomConfigCancel
			return
		}

		err = validateIqResponse(reply)
		if err != nil {
			log.WithError(ErrInformationQueryResponse).Error("An error occurred when trying to read the response from the room configuration rollback request")
			ec <- ErrRoomConfigCancelResponse
			return
		}

		close(ec)
	}()

	return ec
}

func validateIqResponse(reply <-chan data.Stanza) error {
	stanza, ok := <-reply
	if !ok {
		return ErrInformationQueryResponse
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return ErrInformationQueryResponse
	}

	if iq.Type == "error" {
		if iq.Error.MUCConflict != nil {
			return ErrOwnerAffiliationRevokeConflict
		}

		if iq.Error.MUCNotAllowed != nil {
			return ErrNotAllowedKickOccupant
		}

		return ErrUnexpectedResponse
	}

	return nil
}
