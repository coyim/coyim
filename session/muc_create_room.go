package session

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// TODO: We should refactor EVERYWHERE so that for the room bla@service.example.org
// the local part, which is "bla" should be called roomName
// the full thing "bla@service.example.org" should be called roomID

var (
	// ErrInvalidInformationQueryRequest is an invalid information query request error
	ErrInvalidInformationQueryRequest = errors.New("invalid information query request")

	// ErrUnexpectedResponse is an unexpected response from the server error
	ErrUnexpectedResponse = errors.New("received an unexpected response from the server")

	// ErrInformationQueryResponse contains an error received in the information query response
	ErrInformationQueryResponse = errors.New("received an error from the server")
)

func newCreateMUCRoomContext(s *session, ident jid.Bare) *createMUCRoomContext {
	c := &createMUCRoomContext{
		ident:       ident,
		errorResult: make(chan error),
		s:           s,
	}

	return c
}

// TODO: Add a RoomConfigurationQuery for create a Reserved Room
func (s *session) CreateRoom(ident jid.Bare) <-chan error {
	c := newCreateMUCRoomContext(s, ident)
	go c.createRoom()
	return c.errorResult
}

type createMUCRoomContext struct {
	ident jid.Bare
	// TODO: Maybe rename to errorChannel for consistency?
	errorResult chan error
	s           *session
}

func (c *createMUCRoomContext) createRoom() {
	// See XEP-0045 v1.32.0, section: 10.1.1
	err := c.sendMUCPresence()
	if err != nil {
		c.errorResult <- err
		return
	}

	// // See XEP-0045 v1.32.0, section: 10.1.2
	reply, err := c.sendInformationQuery()
	if err != nil {
		c.errorResult <- ErrUnexpectedResponse
		return
	}

	err = c.checkForErrorsInResponse(reply)
	if err != nil {
		c.logWithError(err, "Invalid information query response")
		c.errorResult <- err
		return
	}

	close(c.errorResult)
}

// TODO: I do not like the name of this method. I have really no idea what it means!
func (c *createMUCRoomContext) identity() string {
	return c.ident.String()
}

func (c *createMUCRoomContext) logWithError(e error, m string) {
	c.s.log.WithError(e).Error(m)
}

func (c *createMUCRoomContext) sendMUCPresence() error {
	err := c.s.conn.SendMUCPresence(c.identity())
	if err != nil {
		c.logWithError(err, "An error ocurred while sending a presence for creating an instant room")
		return ErrUnexpectedResponse
	}
	return nil
}

func (c *createMUCRoomContext) newRoomConfiguration() data.MUCRoomConfiguration {
	return data.MUCRoomConfiguration{
		Form: &data.Form{
			Type: "submit",
		},
	}
}

func (c *createMUCRoomContext) sendInformationQuery() (<-chan data.Stanza, error) {
	reply, _, err := c.s.conn.SendIQ(c.identity(), "set", c.newRoomConfiguration())
	if err != nil {
		c.logWithError(err, "An error ocurred while sending the information query for creating an instant room")
		return nil, err
	}
	return reply, nil
}

// TODO: This method lies about what it is doing - it is waiting for a response
// and then checking it for an error. The name should reflect that
// maybe "waitAndCheckResponse"
func (c *createMUCRoomContext) checkForErrorsInResponse(reply <-chan data.Stanza) error {
	stanza, ok := <-reply
	if !ok {
		return ErrInvalidInformationQueryRequest
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return ErrUnexpectedResponse
	}

	if iq.Type == "error" {
		return ErrInformationQueryResponse
	}

	return nil
}
