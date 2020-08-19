package session

import (
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// ErrInvalidInformationQueryRequest is an invalid information query request error
type ErrInvalidInformationQueryRequest struct{}

// ErrUnexpectedResponse is an unexpected response from the server error
type ErrUnexpectedResponse struct{}

// ErrInformationQueryResponse contains an error received in the information query response
type ErrInformationQueryResponse struct {
	Type   string
	Reason string
}

func (e *ErrInvalidInformationQueryRequest) Error() string {
	return "invalid information query request"
}

func (e *ErrUnexpectedResponse) Error() string {
	return "received an unexpected response from the server"
}

func (e *ErrInformationQueryResponse) Error() string {
	return e.Reason
}

func validateStanza(reply <-chan data.Stanza) error {
	stanza, ok := <-reply
	if !ok {
		return &ErrInvalidInformationQueryRequest{}
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return &ErrUnexpectedResponse{}
	}

	if iq.Type == "error" {
		return &ErrInformationQueryResponse{
			Type:   iq.Type,
			Reason: iq.Error.Text,
		}
	}

	return nil
}

func newRoomConfiguration() data.MUCRoomConfiguration {
	return data.MUCRoomConfiguration{
		Form: &data.Form{
			Type: "submit",
		},
	}
}

// Send a presence for creating the room and signals support for MUC
func (s *session) createRoom(roomID jid.Bare, errorResult chan<- error) {
	// See XEP-0045 v1.32.0, section: 10.1.1
	err := s.conn.SendMUCPresence(roomID.String())
	if err != nil {
		s.log.WithError(err).Error("An error ocurred while sending a presence for creating an instant room")
		errorResult <- &ErrUnexpectedResponse{}
		return
	}

	// See XEP-0045 v1.32.0, section: 10.1.2
	reply, _, err := s.conn.SendIQ(roomID.String(), "set", newRoomConfiguration())
	if err != nil {
		s.log.WithError(err).Error("An error ocurred while sending the information query for creating an instant room")
		errorResult <- &ErrUnexpectedResponse{}
		return
	}

	err = validateStanza(reply)
	if err != nil {
		s.log.WithError(err).Error("Invalid information query response")
		errorResult <- err
		return
	}

	close(errorResult)
}

// TODO: Add a RoomConfigurationQuery for create a Reserved Room
func (s *session) CreateRoom(roomID jid.Bare) <-chan error {
	errorResult := make(chan error)
	go s.createRoom(roomID, errorResult)

	return errorResult
}
