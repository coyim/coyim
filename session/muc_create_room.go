package session

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// CreateRoomError represents an error from create room functionality
type CreateRoomError error

var (
	// ErrInvalidInformationQueryRequest is an invalid information query request error
	ErrInvalidInformationQueryRequest CreateRoomError = errors.New("invalid information query request")

	// ErrUnexpectedResponse is an unexpected response from the server error
	ErrUnexpectedResponse CreateRoomError = errors.New("received an unexpected response from the server")

	// ErrInformationQueryResponse contains an error received in the information query response
	ErrInformationQueryResponse CreateRoomError = errors.New("received an error from the server")
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
	ident       jid.Bare
	errorResult chan error
	s           *session
}

// Send a presence for creating the room and signals support for MUC
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

func (c *createMUCRoomContext) identity() string {
	return c.ident.String()
}

func (c *createMUCRoomContext) logWithError(err error, message string) {
	c.s.log.WithError(err).Error(message)
}

func (c *createMUCRoomContext) sendMUCPresence() CreateRoomError {
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

func (c *createMUCRoomContext) checkForErrorsInResponse(reply <-chan data.Stanza) CreateRoomError {
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
