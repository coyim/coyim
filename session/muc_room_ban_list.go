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
func (s *session) GetRoomBanList(roomID jid.Bare) (<-chan muc.RoomBanList, <-chan error) {
	lc := make(chan muc.RoomBanList)
	ec := make(chan error)

	l := s.muc.log.WithField("room", roomID)

	go func() {
		items, err := s.muc.requestRoomBanList(roomID)
		if err != nil {
			l.WithError(err).Error("Server returned a weird result when requesting the room's ban list.")
			ec <- err
			return
		}

		list := muc.RoomBanList{}

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
		return nil, err
	}

	iq := <-stanza

	reply, ok := iq.Value.(*xmppData.ClientIQ)
	if !ok || reply.Type != "result" {
		return nil, errors.New("the client iq reply is not the expected")
	}

	var list xmppData.MUCAdmin
	if err := xml.Unmarshal(reply.Query, &list); err != nil {
		return nil, errors.New("failed to parse room's ban list response")
	}

	return list.Items, nil
}

func newRoomBanListRequestQuery() xmppData.MUCAdmin {
	return xmppData.MUCAdmin{
		Items: []xmppData.MUCItem{
			{Affiliation: data.AffiliationOutcast},
		},
	}
}

// ModifyRoomBanList modifies the ban list for the given room with the changed items
func (s *session) ModifyRoomBanList(roomID jid.Bare, changedItems []*muc.RoomBanListItem) (<-chan bool, <-chan error) {
	rc := make(chan bool)
	ec := make(chan error)

	go func() {
		err := s.muc.modifyRoomBanList(roomID, changedItems)
		if err != nil {
			s.muc.log.WithField("room", roomID).WithError(err).Error("The ban list of the room, can't be updated")
			ec <- err
		} else {
			rc <- true
		}
	}()

	return rc, ec
}

func (m *mucManager) modifyRoomBanList(roomID jid.Bare, items []*muc.RoomBanListItem) error {
	stanza, _, err := m.conn().SendIQ(roomID.String(), "set", newRoomBanListSaveQuery(items))
	if err != nil {
		return err
	}

	iq := <-stanza

	reply, ok := iq.Value.(*xmppData.ClientIQ)
	if !ok || reply.Type != "result" {
		return errors.New("the client iq reply is not the expected")
	}

	return nil
}

func newRoomBanListSaveQuery(items []*muc.RoomBanListItem) xmppData.MUCAdmin {
	list := []xmppData.MUCItem{}
	for _, itm := range items {
		list = append(list, xmppData.MUCItem{
			Jid:         itm.Jid.String(),
			Affiliation: itm.Affiliation.Name(),
			Reason:      itm.Reason,
		})
	}

	return xmppData.MUCAdmin{Items: list}
}

// UpdateOccupantAffiliations modifies the affiliations of a list of occupants in the given room
func (s *session) UpdateOccupantAffiliations(roomID jid.Bare, occupantAffiliations []*muc.RoomOccupantItem) (<-chan bool, <-chan error) {
	rc := make(chan bool)
	ec := make(chan error)

	go func() {
		ok, err := s.muc.updateOccupantAffiliations(roomID, occupantAffiliations)
		if err != nil {
			s.muc.log.WithField("room", roomID).WithError(err).Error("The occupant affiliations couldn't be updated")
			ec <- err
		} else {
			rc <- ok
		}
	}()

	return rc, ec
}

func (m *mucManager) updateOccupantAffiliations(roomID jid.Bare, occupantAffiliations []*muc.RoomOccupantItem) (bool, error) {
	// All xmpp servers should accept the modification of more than one occupant's affiliation
	// in a single request but the current version of Prosody server doesn't support that.
	// So, this is a temporary iteration that was implemented in order to
	// configure more than one occupant affiliation on Prosody servers.
	for _, i := range occupantAffiliations {
		err := m.modifyOccupantAffiliation(roomID, i)
		if err != nil {
			m.log.WithFields(log.Fields{
				"occupant":    i.Jid.String(),
				"affiliation": i.Affiliation.Name(),
				"reason":      i.Reason,
			}).WithError(err).Error("The occupant affiliation couldn't be updated")
			return false, err
		}
	}
	return true, nil
}

func (m *mucManager) modifyOccupantAffiliation(roomID jid.Bare, occupantAffiliations *muc.RoomOccupantItem) error {
	stanza, _, err := m.conn().SendIQ(roomID.String(), "set", newRoomOccupantAffiliationQuery(occupantAffiliations))
	if err != nil {
		return err
	}

	iq := <-stanza

	reply, ok := iq.Value.(*xmppData.ClientIQ)
	if !ok || reply.Type != "result" {
		return errors.New("the client iq reply is not the expected")
	}

	return nil
}

func newRoomOccupantAffiliationQuery(item *muc.RoomOccupantItem) xmppData.MUCAdmin {
	return xmppData.MUCAdmin{
		Items: []xmppData.MUCItem{
			{
				Jid:         item.Jid.String(),
				Affiliation: item.Affiliation.Name(),
				Reason:      item.Reason,
			},
		},
	}
}
