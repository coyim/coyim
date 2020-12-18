package session

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (s *session) SubmitRoomConfigurationForm(roomID jid.Bare, form *muc.RoomConfigForm) (chan bool, chan error) {
	log := log.WithFields(log.Fields{
		"room":  roomID,
		"where": "SubmitRoomConfigurationForm",
	})

	succesChannel := make(chan bool)
	errorChannel := make(chan error)
	go func() {
		reply, _, err := s.conn.SendIQ(roomID.String(), "set", data.MUCRoomConfiguration{
			Form: form.GetFormData(),
		})

		if err != nil {
			log.WithError(err).Error("An error ocured trying to send iq for configuring room")
			errorChannel <- err
			return
		}

		stanza, ok := <-reply
		if !ok {
			log.WithError(errors.New("error reading stanza reply")).Error("An error ocured trying to read the stanza reply")
			errorChannel <- err
			return
		}

		err = validateStanza(stanza)
		if err != nil {
			log.WithError(err).Error("Error in stanza reply")
			errorChannel <- err
			return
		}
		succesChannel <- true

	}()
	return succesChannel, errorChannel
}

func validateStanza(stanza data.Stanza) error {
	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return ErrUnexpectedResponse
	}

	if iq.Type == "error" {
		return ErrInformationQueryResponse
	}

	return nil
}
