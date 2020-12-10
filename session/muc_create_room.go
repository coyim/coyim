package session

import (
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

var (
	// ErrInvalidReserveRoomRequest is an invalid room reservation request error
	ErrInvalidReserveRoomRequest = errors.New("invalid reserve room request")

	// ErrInvalidInformationQueryRequest is an invalid information query request error
	ErrInvalidInformationQueryRequest = errors.New("invalid information query request")

	// ErrUnexpectedResponse is an unexpected response from the server error
	ErrUnexpectedResponse = errors.New("received an unexpected response from the server")

	// ErrInformationQueryResponse contains an error received in the information query response
	ErrInformationQueryResponse = errors.New("received an error from the server")
)

func (s *session) CreateInstantRoom(roomID jid.Bare) <-chan error {
	c := s.newCreateMUCRoomContext(roomID)
	return c.createInstantRoom()
}

func (s *session) CreateReservedRoom(roomID jid.Bare) (<-chan *muc.RoomConfigForm, <-chan error) {
	c := s.newCreateMUCRoomContext(roomID)
	return c.createReservedRoom()
}

type createMUCRoomContext struct {
	roomID            jid.Bare
	configFormChannel chan *muc.RoomConfigForm
	errorChannel      chan error
	conn              xi.Conn
	log               coylog.Logger
}

func (s *session) newCreateMUCRoomContext(roomID jid.Bare) *createMUCRoomContext {
	c := &createMUCRoomContext{
		roomID:       roomID,
		errorChannel: make(chan error),
		conn:         s.conn,
		log:          s.log.WithField("where", "createRoomContext"),
	}

	return c
}

// See XEP-0045 v1.32.0, section: 10.1.1
func (c *createMUCRoomContext) reserveRoom() bool {
	err := c.sendMUCPresence()
	if err != nil {
		c.errorChannel <- err
		return false
	}
	return true
}

func (c *createMUCRoomContext) createRoom(sendIQ func() (<-chan data.Stanza, error), onStanzaReceived func(stanza data.Stanza) (string, error)) {
	if !c.reserveRoom() {
		c.finishWithError(ErrInvalidReserveRoomRequest, "An error occurred while reserving the room")
		return
	}

	reply, err := sendIQ()
	if err != nil {
		c.finishWithError(ErrUnexpectedResponse, "Unexpected information query response")
		return
	}

	stanza, ok := <-reply
	if !ok {
		c.finishWithError(ErrInvalidInformationQueryRequest, "Unexpected information query reply")
		return
	}

	errMessage, err := onStanzaReceived(stanza)
	if err != nil {
		c.finishWithError(err, errMessage)
		return
	}

	close(c.errorChannel)
}

// createInstantRoom will create a room "instantly" accepting the default configuration of the room
// For more information see XEP-0045 v1.32.0, section: 10.1.2
func (c *createMUCRoomContext) createInstantRoom() <-chan error {
	go c.createRoom(c.sendIQForInstantRoom, func(stanza data.Stanza) (string, error) {
		err := c.validateStanzaReceived(stanza)
		if err != nil {
			return "Invalid information query response", err
		}
		return "", nil
	})

	return c.errorChannel
}

// createReservedRoom will reserve a room and request the configuration form for it
func (c *createMUCRoomContext) createReservedRoom() (<-chan *muc.RoomConfigForm, <-chan error) {
	c.configFormChannel = make(chan *muc.RoomConfigForm)

	go c.createRoom(c.sendIQForReservedRoom, func(stanza data.Stanza) (string, error) {
		form, err := c.getConfigFormFromStanza(stanza)
		if err != nil {
			return "Invalid information query response", err
		}

		c.configFormChannel <- form

		return "", nil
	})

	return c.configFormChannel, c.errorChannel
}

func (c *createMUCRoomContext) finishWithError(err error, m string) {
	c.logWithError(err, m)
	c.errorChannel <- err
}

func (c *createMUCRoomContext) logWithError(err error, m string) {
	c.log.WithError(err).Error(m)
}

func (c *createMUCRoomContext) sendMUCPresence() error {
	err := c.conn.SendMUCPresence(c.roomID.String(), &data.MUC{})
	if err != nil {
		c.logWithError(err, "An error ocurred while sending a presence for creating an instant room")
		return ErrUnexpectedResponse
	}

	return nil
}

func (c *createMUCRoomContext) sendInformationQuery(tp string, d interface{}) (<-chan data.Stanza, error) {
	reply, _, err := c.conn.SendIQ(c.roomID.String(), tp, d)
	if err != nil {
		c.logWithError(err, "An error ocurred while sending the information query")
		return nil, err
	}

	return reply, nil
}

func (c *createMUCRoomContext) sendIQForInstantRoom() (<-chan data.Stanza, error) {
	return c.sendInformationQuery("set", c.newRoomConfigurationFormSubmit())
}

func (c *createMUCRoomContext) sendIQForReservedRoom() (<-chan data.Stanza, error) {
	return c.sendInformationQuery("get", c.newRoomConfigurationFormRequest())
}

func (c *createMUCRoomContext) getIQFromStanza(stanza data.Stanza) (*data.ClientIQ, error) {
	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return nil, ErrUnexpectedResponse
	}

	if iq.Type == "error" {
		return nil, ErrInformationQueryResponse
	}

	return iq, nil
}

func (c *createMUCRoomContext) getConfigFormFromStanza(stanza data.Stanza) (*muc.RoomConfigForm, error) {
	iq, err := c.getIQFromStanza(stanza)
	if err != nil {
		return nil, err
	}

	cf, err := c.getConfigFormFromIQResponse(iq)
	if err != nil {
		return nil, err
	}

	return muc.NewRoomConfigRom(cf.Form), nil
}

func (c *createMUCRoomContext) getConfigFormFromIQResponse(iq *data.ClientIQ) (cf *data.MUCRoomConfiguration, err error) {
	err = xml.Unmarshal(iq.Query, &cf)
	return
}

func (c *createMUCRoomContext) validateStanzaReceived(stanza data.Stanza) error {
	_, err := c.getIQFromStanza(stanza)
	return err
}

func (c *createMUCRoomContext) newRoomConfigurationFormSubmit() data.MUCRoomConfiguration {
	return data.MUCRoomConfiguration{
		Form: &data.Form{
			Type: "submit",
		},
	}
}

func (c *createMUCRoomContext) newRoomConfigurationFormRequest() data.MUCRoomConfiguration {
	return data.MUCRoomConfiguration{}
}
