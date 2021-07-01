package session

import (
	"fmt"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

type mucRequestType string

func (rt mucRequestType) String() string {
	return fmt.Sprintf("%s", rt)
}

type informationQueryType string

const (
	informationQueryTypeGet informationQueryType = "get"
	informationQueryTypeSet informationQueryType = "set"
)

func (qt informationQueryType) String() string {
	return fmt.Sprintf("%s", qt)
}

type mucRequest struct {
	roomID       jid.Bare
	conn         xi.Conn
	errorChannel chan error
	onResponse   func(response []byte) error
	log          coylog.Logger
}

func (m *mucManager) newMUCRoomRequest(roomID jid.Bare, requestType mucRequestType, onResponse func(response []byte) error) *mucRequest {
	return &mucRequest{
		roomID:       roomID,
		conn:         m.conn(),
		errorChannel: make(chan error),
		onResponse:   onResponse,
		log: m.log.WithFields(log.Fields{
			"where":       "mucRequest",
			"requestType": requestType.String(),
		}),
	}
}

func (r *mucRequest) get(query interface{}) {
	r.send(informationQueryTypeGet, query)
}

func (r *mucRequest) set(query interface{}) {
	r.send(informationQueryTypeSet, query)
}

func (r *mucRequest) send(queryType informationQueryType, query interface{}) {
	reply, _, err := r.conn.SendIQ(r.roomID.String(), queryType.String(), query)
	if err != nil {
		r.error(ErrUnexpectedResponse)
		return
	}

	stanza, ok := <-reply
	if !ok {
		r.error(ErrInvalidInformationQueryRequest)
		return
	}

	iq, err := r.clientIQFromStanza(stanza)
	if err != nil {
		r.error(err)
		return
	}

	if r.onResponse != nil {
		err = r.onResponse(iq.Query)
		if err != nil {
			r.error(err)
		}
	}
}

func (r *mucRequest) clientIQFromStanza(stanza data.Stanza) (*data.ClientIQ, error) {
	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return nil, ErrUnexpectedResponse
	}

	if iq.Type == "error" {
		return nil, ErrInformationQueryResponse
	}

	return iq, nil
}

func (r *mucRequest) sendMUCPresence() bool {
	if err := r.conn.SendMUCPresence(r.roomID.String(), &data.MUC{}); err != nil {
		r.error(ErrUnexpectedResponse)
		return false
	}
	return true
}

func (r *mucRequest) error(err error) {
	requestError := r.newMUCRoomRequestError(err)
	requestError.logError()

	if r.errorChannel != nil {
		r.errorChannel <- requestError.err
	}
}
