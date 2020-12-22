package session

import (
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

		err = validateIqResponse(reply)
		if err != nil {
			log.WithError(ErrInformationQueryResponse).Error("An error ocured trying to read the stanza reply")
			errorChannel <- err
			return
		}
		succesChannel <- true

	}()
	return succesChannel, errorChannel
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
		return ErrUnexpectedResponse
	}

	return nil
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
			log.WithError(err).Error("An error ocured trying to send iq for canceling room configuration")
			ec <- err
			return
		}

		err = validateIqResponse(reply)
		if err != nil {
			log.WithError(ErrInformationQueryResponse).Error("An error ocured trying to read the stanza reply")
			ec <- err
			return
		}

		close(ec)
	}()

	return ec
}
