package session

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// TODO[OB]-MUC: This is a fairly large method. Would it be possible to break it up into smaller, more helpful bits

func (s *session) createRoom(roomID jid.Bare, errorResult chan<- error) {
	// Send a presence for create the room and signals support for MUC
	// See: 10.1.1 Create room General
	err := s.conn.SendMUCPresence(roomID.String())
	if err != nil {
		errorResult <- err
		return
	}

	// TODO[OB]-MUC: What does this comment mean?
	// TODO: Delete 'roomConf' and get this information from the function.
	// See: 10.1.2 Creating an Instant Room
	// Room Information Query
	roomConf := &data.MUCRoomConfiguration{
		Form: &data.Form{
			Type: "submit",
		},
	}

	reply, _, err := s.conn.SendIQ(roomID.String(), "set", roomConf)
	if err != nil {
		errorResult <- err
		return
	}

	stanza, ok := <-reply
	if !ok {
		errorResult <- errors.New("xmpp: failed to receive response")
		return
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		// TODO[OB]-MUC: These error messages are not exactly correct. It would be better if they said something more useful about what happened
		err = errors.New("xmpp: failed getting the response")
		s.log.WithError(err).Error("failed getting the response when configuring room")
		errorResult <- err
		return
	}

	if iq.Type == "error" {
		err = errors.New("xmpp: error type response getting from Information Query")
		s.log.WithError(err).Error("error stanza information:", iq.Error.Type, iq.Error.Text)
		errorResult <- err
		return
	}

	// TODO[OB]-MUC: I don't like this pattern of sending a nil on the errorResult channel. Much better to close the channel
	errorResult <- nil
}

// TODO: Add a RoomConfigurationQuery for create a Reserved Room
func (s *session) CreateRoom(roomID jid.Bare) <-chan error {
	// TODO[OB]-MUC: I don't think this channel should be buffered
	errorResult := make(chan error, 1)

	go s.createRoom(roomID, errorResult)

	return errorResult
}
