package session

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

var (
	// ErrInvalidInformationQueryRequest is an invalid information query request error
	ErrInvalidInformationQueryRequest = errors.New("invalid information query request")

	// ErrUnexpectedResponse is an unexpected response from the server error
	ErrUnexpectedResponse = errors.New("received an unexpected response from the server")

	// ErrInformationQueryResponse contains an error received in the information query response
	ErrInformationQueryResponse = errors.New("received an error from the server")
)

func newCreateMUCRoomContext(s *session, roomID jid.Bare) *createMUCRoomContext {
	c := &createMUCRoomContext{
		roomID:       roomID,
		errorChannel: make(chan error),
		s:            s,
	}

	return c
}

// TODO: Add a RoomConfigurationQuery for create a Reserved Room
func (s *session) CreateRoom(roomID jid.Bare) <-chan error {
	c := newCreateMUCRoomContext(s, roomID)
	go c.createRoom()
	return c.errorChannel
}

func (s *session) ReserveRoom(roomID jid.Bare) (<-chan *data.MUCRoomConfiguration, <-chan error) {
	c := newCreateMUCRoomContext(s, roomID)
	c.roomConfigFormChannel = make(chan *data.MUCRoomConfiguration)
	go c.reserveRoom()
	return c.roomConfigFormChannel, c.errorChannel
}

type createMUCRoomContext struct {
	roomID                jid.Bare
	roomConfigFormChannel chan *data.MUCRoomConfiguration
	errorChannel          chan error
	s                     *session
}

func (c *createMUCRoomContext) createRoom() {
	// See XEP-0045 v1.32.0, section: 10.1.1
	err := c.sendMUCPresence()
	if err != nil {
		c.errorChannel <- err
		return
	}

	// // See XEP-0045 v1.32.0, section: 10.1.2
	reply, err := c.sendInformationQuery()
	if err != nil {
		c.errorChannel <- ErrUnexpectedResponse
		return
	}

	stanza, ok := <-reply
	if !ok {
		c.errorChannel <- ErrInvalidInformationQueryRequest
		return
	}

	err = c.validateStanzaReceived(stanza)
	if err != nil {
		c.logWithError(err, "Invalid information query response")
		c.errorChannel <- err
		return
	}

	close(c.errorChannel)
}

func (c *createMUCRoomContext) reserveRoom() {
	err := c.sendMUCPresence()
	if err != nil {
		c.errorChannel <- err
		return
	}

	form, err := c.s.conn.SendConfigurationFormRequest(c.roomID.String())
	if err != nil {
		c.logWithError(err, "An error ocurred while sending the information query for creating an reserved room")
		c.errorChannel <- err
		return
	}

	c.roomConfigFormChannel <- form
}

func (c *createMUCRoomContext) logWithError(e error, m string) {
	c.s.log.WithError(e).Error(m)
}

func (c *createMUCRoomContext) sendMUCPresence() error {
	err := c.s.conn.SendMUCPresence(c.roomID.String(), &data.MUC{})
	if err != nil {
		c.logWithError(err, "An error ocurred while sending a presence for creating an instant room")
		return ErrUnexpectedResponse
	}
	return nil
}

func (c *createMUCRoomContext) newRoomConfigurationWithSubmit() data.MUCRoomConfiguration {
	return data.MUCRoomConfiguration{
		Form: &data.Form{
			Type: "submit",
		},
	}
}

func (c *createMUCRoomContext) sendInformationQuery() (<-chan data.Stanza, error) {
	reply, _, err := c.s.conn.SendIQ(c.roomID.String(), "set", c.newRoomConfigurationWithSubmit())
	if err != nil {
		c.logWithError(err, "An error ocurred while sending the information query for creating an instant room")
		return nil, err
	}
	return reply, nil
}

func (c *createMUCRoomContext) validateStanzaReceived(stanza data.Stanza) error {
	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return ErrUnexpectedResponse
	}

	if iq.Type == "error" {
		return ErrInformationQueryResponse
	}

	return nil
}
