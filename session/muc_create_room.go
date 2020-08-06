package session

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) createRoom(roomID jid.Bare, errorResult chan<- error) {
	// Send a presence for create the room and signals support for MUC
	// See: 10.1.1 Create room General
	err := s.conn.SendMUCPresence(roomID.String())
	if err != nil {
		errorResult <- err
		return
	}

	// TODO: Delete 'roomConf' and get this information from the function.
	// See: 10.1.2 Creating an Instant Room
	// Room Information Query
	roomConf := &data.MUCOwnerQuery{
		Form: &data.Form{
			Type: "submit",
		},
	}

	// Send an IQ for create a Instant Room.
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

	errorResult <- nil
}

// TODO: Add a RoomConfigurationQuery for create a Reserved Room
func (s *session) CreateRoom(roomID jid.Bare) <-chan error {
	errorResult := make(chan error, 1)

	go s.createRoom(roomID, errorResult)

	return errorResult
}

//GetChatServices offers the chat services from a xmpp server.
func (s *session) GetChatServices(server jid.Domain) ([]data.DiscoveryItem, error) {
	items, err := s.conn.QueryServiceItems(server.String())
	if err != nil {
		return nil, err
	}

	return items.DiscoveryItems, nil
}
