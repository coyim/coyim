package session

import (
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

var (
	// ErrInvalidReserveRoomRequest is an invalid room reservation request error
	ErrInvalidReserveRoomRequest = errors.New("invalid reserve room request")
)

const (
	mucRequestCreateInstantRoom  mucRequestType = "create-instant-room"
	mucRequestCreateReservedRoom mucRequestType = "create-reserved-room"
)

// CreateInstantRoom will create a room "instantly" accepting the default configuration of the room
// For more information see XEP-0045 v1.32.0, section: 10.1.2
func (s *session) CreateInstantRoom(roomID jid.Bare) (<-chan bool, <-chan error) {
	resultChannel := make(chan bool)

	ctx := s.muc.newCreateMUCRoomContext(roomID, mucRequestCreateInstantRoom, func(data []byte) error {
		resultChannel <- true
		return nil
	})

	go ctx.createRoom(informationQueryTypeSet, data.MUCRoomConfiguration{
		Form: &data.Form{
			Type: "submit",
		},
	})

	return resultChannel, ctx.errorChannel
}

// CreateReservedRoom will reserve a room and request the configuration form for it
func (s *session) CreateReservedRoom(roomID jid.Bare) (<-chan *muc.RoomConfigForm, <-chan error) {
	configFormChannel := make(chan *muc.RoomConfigForm)

	ctx := s.muc.newCreateMUCRoomContext(roomID, mucRequestCreateReservedRoom, func(response []byte) error {
		cf := &data.MUCRoomConfiguration{}
		if err := xml.Unmarshal(response, cf); err != nil {
			return err
		}

		configFormChannel <- muc.NewRoomConfigForm(cf.Form)
		return nil
	})

	go ctx.createRoom(informationQueryTypeGet, data.MUCRoomConfiguration{})

	return configFormChannel, ctx.errorChannel
}

type createMUCRoomContext struct {
	*mucRequest
}

func (m *mucManager) newCreateMUCRoomContext(roomID jid.Bare, requestType mucRequestType, onResponse func([]byte) error) *createMUCRoomContext {
	return &createMUCRoomContext{
		m.newMUCRoomRequest(roomID, requestType, onResponse),
	}
}

func (c *createMUCRoomContext) createRoom(queryType informationQueryType, query interface{}) {
	// See XEP-0045 v1.32.0, section: 10.1.1
	if ok := c.sendMUCPresence(); ok {
		c.send(queryType, query)
	}
}
