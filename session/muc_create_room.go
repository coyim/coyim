package session

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// TODO: Add a RoomConfigurationQuery for create a Reserved Room
func (s *session) CreateRoom(roomID jid.Bare) error {
	// Send a presence for create the room and signals support for MUC
	// See: 10.1.1 Create room General
	err := s.conn.SendMUCPresence(roomID.String())
	if err != nil {
		return err
	}

	// TODO: Delete 'roomConf' and get this information from the function.
	// See: 10.1.2 Creating an Instant Room
	// Room Information Query
	roomConf := &data.RoomConfigurationQuery{
		Form: &data.Form{
			Type: "submit",
		},
	}

	// Send an IQ for create a Instant Room.
	reply, _, err := s.conn.SendIQ(roomID.String(), "set", roomConf)
	if err != nil {
		return err
	}

	stanza, ok := <-reply
	if !ok {
		return errors.New("xmpp: failed to receive response")
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		err = errors.New("xmpp: failed getting the response")
		s.log.WithError(err).Error("failed getting the response when configuring room")
		return err
	}

	if iq.Type == "error" {
		err = errors.New("xmpp: error type response getting from Information Query")
		s.log.WithError(err).Error("error stanza information:", iq.Error.Type, iq.Error.Text)
		return err
	}

	return err
}
