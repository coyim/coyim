package session

import (
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

// GetRoomBanList can be used to request the banned users list from the given room.
func (s *session) GetRoomBanList(roomID jid.Bare) (<-chan []*muc.RoomBanListItem, <-chan error) {
	lc := make(chan []*muc.RoomBanListItem)
	ec := make(chan error)

	l := s.muc.log.WithField("room", roomID)

	go func() {
		items, err := s.muc.requestRoomBanList(roomID)
		if err != nil {
			l.WithError(err).Error("Server returned a weird result when requesting the room's ban list.")
			ec <- err
			return
		}

		list := []*muc.RoomBanListItem{}

		for _, itm := range items {
			affiliation, err := data.AffiliationFromString(itm.Affiliation)
			if err != nil {
				l.WithError(err).Error("Invalid affiliation from the room ban list item")
				continue
			}

			list = append(list, &muc.RoomBanListItem{
				Jid:         jid.Parse(itm.Jid),
				Affiliation: affiliation,
				Reason:      itm.Reason,
			})
		}

		lc <- list
	}()

	return lc, ec
}

func (m *mucManager) requestRoomBanList(roomID jid.Bare) ([]xmppData.MUCItem, error) {
	stanza, _, err := m.conn().SendIQ(roomID.String(), "get", newRoomBanListRequestQuery())
	if err != nil {
		m.log.WithFields(log.Fields{
			"room":  roomID,
			"where": "requestRoomBanList",
		}).WithError(err).Error("Invalid information query response for the room ban list request")
		return nil, err
	}

	iq := <-stanza

	reply, ok := iq.Value.(*xmppData.ClientIQ)
	if !ok || reply.Type != "result" {
		return nil, errors.New("the client iq reply is not the expected")
	}

	var list xmppData.MUCRoomBanListItems
	if err := xml.Unmarshal(reply.Query, &list); err != nil {
		return nil, errors.New("failed to parse room's ban list response")
	}

	return list.Items, nil
}

func newRoomBanListRequestQuery() xmppData.MUCRoomBanListRequestQuery {
	return xmppData.MUCRoomBanListRequestQuery{
		Item: xmppData.MUCItem{
			Affiliation: data.AffiliationOutcast,
		},
	}
}
